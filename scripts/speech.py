# speech.py
# speechrecognition, pyaudio, brew install portaudio
import speech_recognition as sr


class Speech(object):
    def __init__(self):
        self.r = sr.Recognizer()
        self.r.pause_threshold = 0.8
        self.r.phrase_time_limit = 8
        self.r.operation_timeout = 10
        self.r.dynamic_energy_threshold = True
        # self.r.dynamic_energy_adjustment_ratio = 1.25

    def google_speech_recognition(self, recognizer, audio):
        speech = None
        reason = None

        try:
            speech = recognizer.recognize_google(audio)
        except sr.UnknownValueError:
            print("Google Speech Recognition could not understand audio")
            reason = "could-not-translate"
        except sr.RequestError as e:
            print("Could not request results from Google Speech Recognition service; {0}".format(e))
            reason = "api-error"
        except:
            print("Google Speech timed out")
            reason = "api-timeout"

        return speech, reason

    def listen_for_audio(self, timeout):
        # obtain audio from the microphone
        m = sr.Microphone(chunk_size=4096)
        with m as source:
            self.r.adjust_for_ambient_noise(source, duration=1)

            try:
                print("I'm listening")
                audio = self.r.listen(source, timeout=timeout, phrase_time_limit=10)
            except sr.WaitTimeoutError:
                print("Timed out.")
                return None, None

        # self.__debugger_microphone(enable=False)
        print("Found audio")
        return self.r, audio
