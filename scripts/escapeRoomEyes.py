import sys
import cv2 as cv
import requests
from time import sleep

cap = cv.VideoCapture(0)

frontal_face = cv.CascadeClassifier('models/facial_recognition_model.xml')
profile_face = cv.CascadeClassifier('models/haarcascade_profileface.xml')

lastCount = 0

game_id = "eWpMRNLMxiDx47yqg"


def send_face_count(current_count, previous_count):
    """This sends the face count to the effects server"""
    print("Sending: " + str(current_count))
    r = requests.post("http://localhost:8080/faces",
                      json={"_id": game_id, "previousCount": previous_count, "currentCount": current_count})
    print("Status: " + str(r.status_code))
    return


if len(sys.argv) > 1:
    game_id = sys.argv[1]
    print("Game Id set to " + game_id)
else:
    print("Using default game_id")

# Initialize and test connection
send_face_count(lastCount, 0)

while True:
    # Capture frame-by-frame
    ret, frame = cap.read()

    # Our operations on the frame come here
    gray = cv.cvtColor(frame, cv.COLOR_BGR2GRAY)

    # Find the faces in our image
    faces = frontal_face.detectMultiScale(gray,
                                          scaleFactor=1.1,
                                          minNeighbors=5,
                                          minSize=(30, 30)
                                          )

    profiles = profile_face.detectMultiScale(gray,
                                             scaleFactor=1.1,
                                             minNeighbors=5,
                                             minSize=(30, 30)
                                             )

    faceCount = len(faces) + len(profiles)
    print('Faces:: last: ' + str(lastCount) + ' current: ' + str(faceCount))

    if faceCount != lastCount:
        try:
            send_face_count(faceCount, lastCount)
            lastCount = faceCount
        except:
            print("Failed to send!")

    sleep(0.25)

    if cv.waitKey(1) & 0xFF == ord('q'):
        break

# Release when done
cap.release()
cv.destroyAllWindows()
