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
            self.game_id = args[1]
        else:
            self.game_id = "eWpMRNLMxiDx47yqg"

        print("Game Id set to " + self.game_id)

        self.lastCount = 0
        self.sleepTime = 0.25
        self.speech = Speech(launch_phrase="hey you", debugger_enabled=False)

    def start(self):
        """Main loop"""

        # Initialize and test connection
        self.send_face_count(self.lastCount, 0)

        while True:
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
            print('Faces:: last: ' + str(self.lastCount) + ' current: ' + str(face_count))

            if face_count != self.lastCount:
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
        print("Ready to listen for speech...")
        recognizer, audio = self.speech.listen_for_audio(3)

        if audio is not None:
            # received audio data, now we'll recognize it using Google Speech Recognition
            print("Getting speech...")
            speech = self.speech.google_speech_recognition(recognizer, audio)

            if speech is not None:
                print("Speech: " + speech)

    def send_face_count(self, current_count, previous_count):
        """This sends the face count to the effects server"""

        print("Sending: " + str(current_count))
        r = requests.post("http://localhost:8080/faces",
                          json={
                              "_id": self.game_id,
                              "previousCount": previous_count,
                              "currentCount": current_count
                          })
        print("Status: " + str(r.status_code))
        return


if __name__ == "__main__":
    bot = EscapeRoomEyes(sys.argv)

    try:
        bot.start()
    except KeyboardInterrupt:
        # Release when done
        cap.release()
        cv.destroyAllWindows()
