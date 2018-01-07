import sys
import RPi.GPIO as GPIO

RoAPin = 17  # dt
RoBPin = 18  # clk
EndPoint = "/hours"

handPosition = 0

flag = 0
Last_RoB_Status = 0
Current_RoB_Status = 0


def setup():
    global RoAPin, RoBPin, EndPoint
    RoAPin = int(sys.argv[1])
    RoBPin = int(sys.argv[2])
    EndPoint = sys.argv[3]

    GPIO.setmode(GPIO.BCM)
    GPIO.setup(RoAPin, GPIO.IN)
    GPIO.setup(RoBPin, GPIO.IN)


def lockHandPosition(pos):
    if pos == 24:
        return 0
    elif pos == -1:
        return 23

    return pos


def rotaryDeal():
    global flag
    global Last_RoB_Status
    global Current_RoB_Status
    global handPosition
    Last_RoB_Status = GPIO.input(RoBPin)
    while not GPIO.input(RoAPin):
        Current_RoB_Status = GPIO.input(RoBPin)
        flag = 1
    if flag == 1:
        flag = 0
        if (Last_RoB_Status == 0) and (Current_RoB_Status == 1):
            handPosition = lockHandPosition(handPosition + 1)
            print('handPosition = %d' % handPosition)
        if (Last_RoB_Status == 1) and (Current_RoB_Status == 0):
            handPosition = lockHandPosition(handPosition - 1)
            print('handPosition = %d' % handPosition)


def loop():
    while True:
        rotaryDeal()


def destroy():
    GPIO.cleanup()  # Release resource


if __name__ == '__main__':  # Program start from here
    if len(sys.argv) != 4:
        print("Usage: clockEncoder.py <PIN_A> <PIN_B> <endpoint>")
    else:
        setup()
        try:
            loop()
        except KeyboardInterrupt:
            destroy()
