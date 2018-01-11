# Install OpenCV
http://milq.github.io/install-opencv-ubuntu-debian/

# Install audio
```
sudo apt-get install python-pyaudio python3-pyaudio
sudo -H pip uninstall PyAudio
sudo -H pip install PyAudio
sudo -H pip install SpeechRecognition
```

# Clock Commands
```
sudo python clockEncoder.py  17 18 "http://192.168.86.50:8080/hours"
sudo python clockEncoder.py 5 6 "http://192.168.86.50:8080/minutes"
```