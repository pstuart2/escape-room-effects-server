import sys
import threading
from time import sleep

import cv2 as cv
import requests

from speech import Speech

cap = cv.VideoCapture(0)

frontal_face = cv.CascadeClassifier('models/haarcascade_frontalface_default.xml')
profile_face = cv.CascadeClassifier('models/haarcascade_profileface.xml')

server = "http://localhost:8080"
face_count = 0


def count_faces():
    global face_count
    last_count = 0

    while True:
        ret, frame = cap.read()

        # Our operations on the frame come here
        gray = cv.cvtColor(frame, cv.COLOR_BGR2GRAY)

        # Find the faces in our image
        faces = frontal_face.detectMultiScale(gray,
                                              scaleFactor=1.5,
                                              minNeighbors=5,
                                              minSize=(30, 30)
                                              )

        face_count = len(faces)

        if face_count != last_count:
            print('Faces:: last: ' + str(last_count) + ' current: ' + str(face_count))
            try:
                send_face_count(face_count, last_count)
                last_count = face_count
            except:
                print("Failed to send!")


def send_face_count(current_count, previous_count):
    """This sends the face count to the effects server"""

    print("Sending: " + str(current_count))
    r = requests.post(server + "/faces",
                      json={
                          "previousCount": previous_count,
                          "currentCount": current_count
                      })
    if r.status_code != 200:
        print("send_face_count Status: " + str(r.status_code))


class EscapeRoomEyes(object):
    def __init__(self):
        self.speech = Speech()

    def start(self):
        """Main loop"""

        while True:
            if face_count == 0:
                sleep(0.25)
                continue

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

    def send_command(self, command, text=""):
        """This sends the speech to the effects server for processing"""
        print("Sending command: " + command + ", text: [" + text + "]")

        r = requests.post(server + "/command",
                          json={
                              "command": command,
                              "text": text.lower(),
                          })
        if r.status_code != 200:
            print("send_command Status: " + str(r.status_code))


if __name__ == "__main__":
    if len(sys.argv) > 1:
        server = sys.argv[1]

    # Initialize and test connection
    send_face_count(0, 0)

    bot = EscapeRoomEyes()

    d = threading.Thread(name='count_faces', target=count_faces)
    d.setDaemon(True)

    try:
        d.start()
        bot.start()
    except KeyboardInterrupt:
        # Release when done
        cap.release()
        cv.destroyAllWindows()
