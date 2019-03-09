package processor

import (
  "log"
  "image"
  . "github.com/3d0c/gmf"
  "gocv.io/x/gocv"
)

const MinimumArea = 3000
const WidthPadding = .05
const HeightPadding = .05

type Motion struct {
  MotionDetected bool
  FramePktPts int64
  MotionAreas []image.Rectangle
}

func setFrame(frame *Frame, dstMat *gocv.Mat, srcCodecCtx *CodecCtx, timeBase AVR) error {
  ret, err := frameToMat(frame, srcCodecCtx, timeBase)
  if (err != nil) {
    return err
  }

  matToBW(ret)
  ret.CopyTo(dstMat)
  ret.Close()

  return nil
}

func DetectMotion(frames <-chan *Frame, doneFrames chan<- *Frame, results chan<- *Motion, srcCodecCtx *CodecCtx, timeBase AVR) {

  prevFrame := gocv.NewMat()
  defer prevFrame.Close()

  curFrame  := gocv.NewMat()
  defer curFrame.Close()

  height := srcCodecCtx.Height()
  width  := srcCodecCtx.Width()

  for frame := range frames {

    if prevFrame.Empty() {
      err := setFrame(frame, &prevFrame, srcCodecCtx, timeBase)
      if (err != nil) {
        log.Print("Could set frame to MAT")
      }

      doneFrames <- frame
      continue
    } else {
      err := setFrame(frame, &curFrame, srcCodecCtx, timeBase)
      if (err != nil) {
        log.Print("Could set frame to MAT")
      }
    }

    motion := new(Motion)
    motion.MotionDetected = false
    motion.FramePktPts = frame.PktPts()
    doneFrames <- frame

    frameDelta := gocv.NewMat()
    thresh     := gocv.NewMat()

    gocv.AbsDiff(prevFrame, curFrame, &frameDelta)
    gocv.Threshold(frameDelta, &thresh, 50, 255, gocv.ThresholdBinary)

    kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
    gocv.Dilate(thresh, &thresh, kernel)

    contours := gocv.FindContours(thresh, gocv.RetrievalExternal, gocv.ChainApproxSimple)
    for _, c := range contours {
      area := gocv.ContourArea(c)
      if area < MinimumArea {
        continue
      }

      motion.MotionDetected = true

      rect := gocv.BoundingRect(c)
      x := rect.Min.X
      y := rect.Min.Y
      w := rect.Size().X
      h := rect.Size().Y

      //Apply padding of the motion area
      widthPadding := int(float64(width) * WidthPadding)
      heightPadding := int(float64(height) * HeightPadding)
      x = x-(widthPadding)
      w = w+(widthPadding*2)
      y = y-(heightPadding)
      h = h+(heightPadding*2)

      rectangle := image.Rect(x, y, x+w, y+h)

      motion.MotionAreas = append(motion.MotionAreas, rectangle)
    }

    curFrame.CopyTo(&prevFrame)

    //Return our results to the channel
    results <- motion

    frameDelta.Close()
    thresh.Close()
    kernel.Close()

  }
}
