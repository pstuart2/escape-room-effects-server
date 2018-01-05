import numpy as np
import cv2 as cv

cap = cv.VideoCapture(0)

# face_cascade = cv.CascadeClassifier('models/facial_recognition_model.xml')
frontal_face = cv.CascadeClassifier('models/haarcascade_frontalface_default.xml')
profile_face = cv.CascadeClassifier('models/haarcascade_profileface.xml')

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

    for (x, y, w, h) in faces:
        cv.rectangle(frame, (x, y), (x + w, y + h), (255, 0, 0), 2)

    profiles = profile_face.detectMultiScale(gray,
                                             scaleFactor=1.5,
                                             minNeighbors=5,
                                             minSize=(30, 30)
                                             )

    for (x, y, w, h) in profiles:
        cv.rectangle(frame, (x, y), (x + w, y + h), (0, 255, 0), 2)

    # Display the resulting frame
    cv.imshow('frame', frame)

    if cv.waitKey(1) & 0xFF == ord('q'):
        break

# Release when done
cap.release()
cv.destroyAllWindows()
