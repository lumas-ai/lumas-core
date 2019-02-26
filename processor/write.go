package processor

import (
  "fmt"
  "log"
  "strconv"
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
  dstFileName := strconv.FormatInt(s.Id, 10) + ".ts"

  outputCtx := assert(NewOutputCtxWithFormatName(dstFileName, "mpegts")).(*FmtCtx)
  defer outputCtx.Close()
  outputCtx.SetStartTime(0)
  addStreams(inputCtx, outputCtx)
  outputCtx.Dump()

  if err := outputCtx.WriteHeader(); err != nil {
    fmt.Println("error")
  }
  defer outputCtx.WriteTrailer()

  for packet := range packets {
    //Write the packet to a file
    if err := outputCtx.WritePacket(packet); err != nil {
      fmt.Println("Could not write packet")
    }

    donePackets <- packet
  }
}
