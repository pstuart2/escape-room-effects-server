import sys
import cv2 as cv
import requests
from time import sleep
from speech import Speech

cap = cv.VideoCapture(0)

frontal_face = cv.CascadeClassifier('models/haarcascade_frontalface_default.xml')
profile_face = cv.CascadeClassifier('models/haarcascade_profileface.xml')


class EscapeRoomEyes(object):
    def __init__(self, args):
        if len(args) > 1:
            self.server = args[1]
        else:
            self.server = "http://localhost:8080"

        print("Server set to " + self.server)

        self.lastCount = 0
        self.sleepTime = 0.25
        self.speech = Speech()

    def start(self):
        """Main loop"""

        # Initialize and test connection
        self.send_face_count(self.lastCount, 0)

        while True:
            #print("-----------------------------------")
            # Capture frame-by-frame
            ret, frame = cap.read()

            # Our operations on the frame come here
            gray = cv.cvtColor(frame, cv.COLOR_BGR2GRAY)

            # Find the faces in our image
            faces = frontal_face.detectMultiScale(gray,
                                                  scaleFactor=1.5,
                                                  minNeighbors=5,
                                                  minSize=(30, 30)
                                                  )

            profiles = profile_face.detectMultiScale(gray,
                                                     scaleFactor=1.5,
                                                     minNeighbors=5,
                                                     minSize=(30, 30)
                                                     )

            face_count = len(faces) + len(profiles)

            if face_count != self.lastCount:
                print('Faces:: last: ' + str(self.lastCount) + ' current: ' + str(face_count))
                try:
                    self.send_face_count(face_count, self.lastCount)
                    self.lastCount = face_count
                except:
                    print("Failed to send!")

            if face_count == 0:
                sleep(self.sleepTime)
                continue

            # recognizer, audio = self.speech.listen_for_audio(1)
            # if self.speech.is_call_to_action(recognizer, audio):
            self.get_speech()

    def get_speech(self):
        self.send_command(":listening")
        recognizer, audio = self.speech.listen_for_audio(3)

        if audio is not None:
            # received audio data, now we'll recognize it using Google Speech Recognition
            self.send_command(":getting-speech")
            speech, reason = self.speech.google_speech_recognition(recognizer, audio)

            if speech is not None:
                self.send_command(":speech", speech)
            else:
                self.send_command(":stopped", reason)

            sleep(3)
        else:
            self.send_command(":stopped", "no-audio")

    def send_face_count(self, current_count, previous_count):
        """This sends the face count to the effects server"""

        print("Sending: " + str(current_count))
        r = requests.post(self.server + "/faces",
                          json={
                              "previousCount": previous_count,
                              "currentCount": current_count
                          })
        if r.status_code != 200:
            print("send_face_count Status: " + str(r.status_code))

    def send_command(self, command, text = ""):
        """This sends the speech to the effects server for processing"""
        print("Sending command: " + command + ", text: [" + text + "]")

        r = requests.post(self.server + "/command",
                          json={
                              "command": command,
                              "text": text.lower(),
                          })
        if r.status_code != 200:
            print("send_command Status: " + str(r.status_code))


if __name__ == "__main__":
    bot = EscapeRoomEyes(sys.argv)

    try:
        bot.start()
    except KeyboardInterrupt:
        # Release when done
        cap.release()
        cv.destroyAllWindows()
