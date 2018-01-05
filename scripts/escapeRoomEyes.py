import sys
import cv2 as cv
import requests
from time import sleep

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

            sleep(self.sleepTime)

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
    bot.start()

    # Release when done
    cap.release()
    cv.destroyAllWindows()
