package processor

import (
  "fmt"
  "log"
  . "github.com/3d0c/gmf"
)

func addStreams(inputCtx *FmtCtx, outputCtx *FmtCtx) {
  for i := 0; i < inputCtx.StreamsCnt(); i++ {
    srcStream, err := inputCtx.GetStream(i)
    if err != nil {
      log.Print("Could not add stream to to output codec context: " + err.Error())
    }

    outputCtx.AddStreamWithCodeCtx(srcStream.CodecCtx())
  }
}

func (s *Camera) WriteFile(packets <-chan *Packet, donePackets chan<- *Packet, inputCtx *FmtCtx) {
  dstFileName := fmt.Sprintf("/videos/%d.ts", s.Id)

  outputCtx := assert(NewOutputCtxWithFormatName(dstFileName, "mpegts")).(*FmtCtx)
  defer outputCtx.Close()
  outputCtx.SetStartTime(0)
  addStreams(inputCtx, outputCtx)
  outputCtx.Dump()

  if err := outputCtx.WriteHeader(); err != nil {
    fmt.Println(fmt.Sprintf("Could not write output file for camera: %d", s.Id))
    fmt.Println(err.Error())
  }
  defer outputCtx.WriteTrailer()

  for packet := range packets {
    //Write the packet to a file
    if err := outputCtx.WritePacket(packet); err != nil {
      fmt.Println(fmt.Sprintf("[%d] Could not write packet", s.Id))
    }

    donePackets <- packet
  }
}
