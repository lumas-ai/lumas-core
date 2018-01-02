from models import object_detection
from threading import Thread
from flask import Flask
import websocket
import thread
import time
import cv2
import sys
import os
import json
import os
import socket
import sys

# Global variables shared between threads
frame = None
ret = None
process_feed = False
camera_open = False
timer = None

# Create the API endpoints
app = Flask(__name__)

def timer():
    global process_feed
    global timer

    seconds = time.time() - timer

    while seconds <= 60:
        sys.stderr.write("Elapsed seconds are: " + seconds)
        time.sleep(1)
        seconds = time.time() - timer

    process_feed = False
    timer = None
    sys.exit(0)

@app.route('/start')
def start():
    global timer
    global process_feed
    sys.stderr.write('Starting feed analysis\n')

    # There must not be threads running if timer == None
    if timer == None:
        timer = time.time()
        timerthread = Thread(target=timer)
        timerthread.start()

    process_feed = True
    return "OK"

@app.route('/stop')
def stop():
    return "OK"

def tf():
    sys.stderr.write("TF thread running\n")

    # Load the TensorFlow models into memory
    base_path = os.path.dirname(os.path.abspath(__file__))
    model_path = base_path + '/faster_rcnn_inception_v2_coco_2017_11_08'
    net = object_detection.Net(graph_fp='%s/frozen_inference_graph.pb' % model_path,
        labels_fp='data/label.pbtxt',
        num_classes=90,
        threshold=0.6)

    while True:
        if not process_feed:
            time.sleep(0.1)
        else:
            if camera_open:
                if ret:
                    resize_frame = cv2.resize(frame, (720, 480))
                    results = net.predict(img=resize_frame, display_img=resize_frame)
                    json_string = json.dumps(results)

                    ws.send(json_string)
        time.sleep(0.1)

    sys.stderr.write("TF thread closing\n")

def camera():
    global camera_open
    global frame
    global ret
    cap = None

    while True:
        if process_feed:
            video_location = os.getenv('STREAM_URL')

            if cap == None:
                try:
                    cap = cv2.VideoCapture(video_location)
                except:
                    sys.stderr.write("Could not open camera feed.\n")
                    camera_open = False

            ret, frame = cap.read()
            camera_open = True
        else:
            if cap:
                cap.release()
                cap = None
            camera_open = False
            time.sleep(0.1)

    sys.stderr.write("Camera thread closing")

def socket():
    global ws

    while True:
        try:
            websocket.enableTrace(False)
            ws = websocket.WebSocketApp("ws://monitor:8080/")
            ws.run_forever()
        except Exception as err:
            sys.stderr.write("Could not connect to websocket: " + str(err) + ". Trying again in 1 second\n")
            time.sleep(1)

def webapp():
    app.run(host='0.0.0.0')

if __name__ == '__main__':
    sys.stderr.write("Starting threads\n")
    camerathread = Thread(target=camera)
    camerathread.start()

    tfthread = Thread(target=tf)
    tfthread.start()

    webthread = Thread(target=webapp)
    webthread.start()
    wsthread = Thread(target=socket)
    wsthread.start()

    wsthread.join()
    tfthread.join()
    camerathread.join()
    webthread.join()

    ws.close()
