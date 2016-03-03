package main

import "flag"
import "fmt"

func main() {

  // Basic flag declarations are available for string,
  // integer, and boolean options. Here we declare a
  // string flag `word` with a default value `"foo"`
  // and a short description. This `flag.String` function
  // returns a string pointer (not a string value);
  // we'll see how to use this pointer below.
  fileLocationPtr := flag.String("file", ".nelson.yml", "path to the nelson YAML file")
  hostPtr := flag.String("host", "nelson.oncue.verizon.net", "host name where the nelson service is running")

  // do the side effects
  flag.Parse()

  // Here we'll just dump out the parsed options and
  // any trailing positional arguments. Note that we
  // need to dereference the pointers with e.g. `*wordPtr`
  // to get the actual option values.
  fmt.Println("file:", *fileLocationPtr)
  fmt.Println("host:", *hostPtr)
}
