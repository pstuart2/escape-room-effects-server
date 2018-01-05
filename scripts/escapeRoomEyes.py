import cv2 as cv
import requests
from time import sleep

cap = cv.VideoCapture(0)

frontal_face = cv.CascadeClassifier('models/facial_recognition_model.xml')
profile_face = cv.CascadeClassifier('models/haarcascade_profileface.xml')

lastCount = 0


def send_face_count(count):
    """This sends the face count to the effects server"""
    print("Sending: " + str(count))
    r = requests.post("http://localhost:8080/faces", json={"count": count})
    print("Status: " + str(r.status_code))
    return


# Initialize and test connection
send_face_count(lastCount)

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
            send_face_count(faceCount)
            lastCount = faceCount
        except:
            print("Failed to send!")

    sleep(0.25)

    if cv.waitKey(1) & 0xFF == ord('q'):
        break

# Release when done
cap.release()
cv.destroyAllWindows()
