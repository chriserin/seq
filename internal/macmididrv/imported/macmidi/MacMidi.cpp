/**********************************************************************/
/*! \class RtMidi
    \brief An abstract base class for realtime MIDI input/output.

    This class implements some common functionality for the realtime
    MIDI input/output subclasses RtMidiIn and RtMidiOut.

    RtMidi GitHub site: https://github.com/thestk/rtmidi
    RtMidi WWW site: http://www.music.mcgill.ca/~gary/rtmidi/

    RtMidi: realtime MIDI i/o C++ classes
    Copyright (c) 2003-2023 Gary P. Scavone

    Permission is hereby granted, free of charge, to any person
    obtaining a copy of this software and associated documentation files
    (the "Software"), to deal in the Software without restriction,
    including without limitation the rights to use, copy, modify, merge,
    publish, distribute, sublicense, and/or sell copies of the Software,
    and to permit persons to whom the Software is furnished to do so,
    subject to the following conditions:

    The above copyright notice and this permission notice shall be
    included in all copies or substantial portions of the Software.

    Any person wishing to distribute modifications to the Software is
    asked to send the modifications to the original developer so that
    they can be incorporated into the canonical version.  This is,
    however, not a binding provision of this license.

    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
    EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
    MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
    IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR
    ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF
    CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
    WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
/**********************************************************************/

#include "MacMidi.h"
#include <sstream>
#if defined(__APPLE__)
#include <TargetConditionals.h>
#endif

// Default for Windows is to add an identifier to the port names; this
// flag can be defined (e.g. in your project file) to disable this behaviour.
// #define RTMIDI_DO_NOT_ENSURE_UNIQUE_PORTNAMES

// Default for Windows UWP is to enable a workaround to fix BLE-MIDI IN ports'
// wrong timestamps that occur at least in Windows 10 21H2;
// this flag can be defined (e.g. in your project file)
// to disable this behavior.
// #define RTMIDI_DO_NOT_ENABLE_WORKAROUND_UWP_WRONG_TIMESTAMPS

// **************************************************************** //
//
// MidiInApi and MidiOutApi subclass prototypes.
//
// **************************************************************** //

#if !defined(__MACOSX_CORE__)
#define __MACMIDI_DUMMY__
#endif

#include <CoreMIDI/CoreMIDI.h>

class MidiInCore : public MidiInApi {
public:
  MidiInCore(const std::string &clientName, unsigned int queueSizeLimit);
  ~MidiInCore(void);
  MacMidi::Api getCurrentApi(void) { return MacMidi::MACOSX_CORE; };
  void openPort(unsigned int portNumber, const std::string &portName);
  void openVirtualPort(const std::string &portName);
  void closePort(void);
  void setClientName(const std::string &clientName);
  void setPortName(const std::string &portName);
  unsigned int getPortCount(void);
  std::string getPortName(unsigned int portNumber);

protected:
  MIDIClientRef
  getCoreMidiClientSingleton(const std::string &clientName) throw();
  void initialize(const std::string &clientName);
};

class MidiOutCore : public MidiOutApi {
public:
  MidiOutCore(const std::string &clientName);
  ~MidiOutCore(void);
  MacMidi::Api getCurrentApi(void) { return MacMidi::MACOSX_CORE; };
  void openPort(unsigned int portNumber, const std::string &portName);
  void openVirtualPort(const std::string &portName);
  void closePort(void);
  void setClientName(const std::string &clientName);
  void setPortName(const std::string &portName);
  unsigned int getPortCount(void);
  std::string getPortName(unsigned int portNumber);
  void sendMessage(const unsigned char *message, size_t size);

protected:
  MIDIClientRef
  getCoreMidiClientSingleton(const std::string &clientName) throw();
  void initialize(const std::string &clientName);
};

#if defined(__MACMIDI_DUMMY__)

class MidiInDummy : public MidiInApi {
public:
  MidiInDummy(const std::string & /*clientName*/, unsigned int queueSizeLimit)
      : MidiInApi(queueSizeLimit) {
    errorString_ = "MidiInDummy: This class provides no functionality.";
    error(MacMidiError::WARNING, errorString_);
  }
  MacMidi::Api getCurrentApi(void) { return MacMidi::MACMIDI_DUMMY; }
  void openPort(unsigned int /*portNumber*/, const std::string & /*portName*/) {
  }
  void openVirtualPort(const std::string & /*portName*/) {}
  void closePort(void) {}
  void setClientName(const std::string & /*clientName*/) {};
  void setPortName(const std::string & /*portName*/) {};
  unsigned int getPortCount(void) { return 0; }
  std::string getPortName(unsigned int /*portNumber*/) { return ""; }

protected:
  void initialize(const std::string & /*clientName*/) {}
};

class MidiOutDummy : public MidiOutApi {
public:
  MidiOutDummy(const std::string & /*clientName*/) {
    errorString_ = "MidiOutDummy: This class provides no functionality.";
    error(MacMidiError::WARNING, errorString_);
  }
  MacMidi::Api getCurrentApi(void) { return MacMidi::MACMIDI_DUMMY; }
  void openPort(unsigned int /*portNumber*/, const std::string & /*portName*/) {
  }
  void openVirtualPort(const std::string & /*portName*/) {}
  void closePort(void) {}
  void setClientName(const std::string & /*clientName*/) {};
  void setPortName(const std::string & /*portName*/) {};
  unsigned int getPortCount(void) { return 0; }
  std::string getPortName(unsigned int /*portNumber*/) { return ""; }
  void sendMessage(const unsigned char * /*message*/, size_t /*size*/) {}

protected:
  void initialize(const std::string & /*clientName*/) {}
};

#endif

//*********************************************************************//
//  MacMidi Definitions
//*********************************************************************//

MacMidi ::MacMidi() : macapi_(0) {}

MacMidi ::~MacMidi() {
  delete macapi_;
  macapi_ = 0;
}

MacMidi::MacMidi(MacMidi &&other) noexcept {
  macapi_ = other.macapi_;
  other.macapi_ = nullptr;
}

std::string MacMidi ::getVersion(void) throw() {
  return std::string(MACMIDI_VERSION);
}

// Define API names and display names.
// Must be in same order as API enum.
extern "C" {
const char *macmidi_api_names[][2] = {
    {"unspecified", "Unknown"},
    {"core", "CoreMidi"},
    {"alsa", "ALSA"},
    {"jack", "Jack"},
    {"winmm", "Windows MultiMedia"},
    {"dummy", "Dummy"},
    {"web", "Web MIDI API"},
    {"winuwp", "Windows UWP"},
    {"amidi", "Android MIDI API"},
};
const unsigned int macmidi_num_api_names =
    sizeof(macmidi_api_names) / sizeof(macmidi_api_names[0]);

// The order here will control the order of MacMidi's API search in
// the constructor.
extern "C" const MacMidi::Api macmidi_compiled_apis[] = {
#if defined(__MACOSX_CORE__)
    MacMidi::MACOSX_CORE,
#endif
#if defined(__LINUX_ALSA__)
    MacMidi::LINUX_ALSA,
#endif
#if defined(__UNIX_JACK__)
    MacMidi::UNIX_JACK,
#endif
#if defined(__WINDOWS_MM__)
    MacMidi::WINDOWS_MM,
#endif
#if defined(__WINDOWS_UWP__)
    MacMidi::WINDOWS_UWP,
#endif
#if defined(__WEB_MIDI_API__)
    MacMidi::WEB_MIDI_API,
#endif
#if defined(__WEB_MIDI_API__)
    MacMidi::WEB_MIDI_API,
#endif
#if defined(__AMIDI__)
    MacMidi::ANDROID_AMIDI,
#endif
    MacMidi::UNSPECIFIED,
};
extern "C" const unsigned int macmidi_num_compiled_apis =
    sizeof(macmidi_compiled_apis) / sizeof(macmidi_compiled_apis[0]) - 1;
}

// This is a compile-time check that macmidi_num_api_names == MacMidi::NUM_APIS.
// If the build breaks here, check that they match.
template <bool b> class StaticAssert {
private:
  StaticAssert() {}
};
template <> class StaticAssert<true> {
public:
  StaticAssert() {}
};
class StaticAssertions {
  StaticAssertions() {
    StaticAssert<macmidi_num_api_names == MacMidi::NUM_APIS>();
  }
};

void MacMidi ::getCompiledApi(std::vector<MacMidi::Api> &apis) throw() {
  apis = std::vector<MacMidi::Api>(
      macmidi_compiled_apis, macmidi_compiled_apis + macmidi_num_compiled_apis);
}

std::string MacMidi ::getApiName(MacMidi::Api api) {
  if (api < MacMidi::UNSPECIFIED || api >= MacMidi::NUM_APIS)
    return "";
  return macmidi_api_names[api][0];
}

std::string MacMidi ::getApiDisplayName(MacMidi::Api api) {
  if (api < MacMidi::UNSPECIFIED || api >= MacMidi::NUM_APIS)
    return "Unknown";
  return macmidi_api_names[api][1];
}

MacMidi::Api MacMidi ::getCompiledApiByName(const std::string &name) {
  unsigned int i = 0;
  for (i = 0; i < macmidi_num_compiled_apis; ++i)
    if (name == macmidi_api_names[macmidi_compiled_apis[i]][0])
      return macmidi_compiled_apis[i];
  return MacMidi::UNSPECIFIED;
}

void MacMidi ::setClientName(const std::string &clientName) {
  macapi_->setClientName(clientName);
}

void MacMidi ::setPortName(const std::string &portName) {
  macapi_->setPortName(portName);
}

//*********************************************************************//
//  MacMidiIn Definitions
//*********************************************************************//

void MacMidiIn ::openMidiApi(MacMidi::Api api, const std::string &clientName,
                             unsigned int queueSizeLimit) {
  delete macapi_;
  macapi_ = 0;

#if defined(__UNIX_JACK__)
  if (api == UNIX_JACK)
    macapi_ = new MidiInJack(clientName, queueSizeLimit);
#endif
#if defined(__LINUX_ALSA__)
  if (api == LINUX_ALSA)
    macapi_ = new MidiInAlsa(clientName, queueSizeLimit);
#endif
#if defined(__WINDOWS_MM__)
  if (api == WINDOWS_MM)
    macapi_ = new MidiInWinMM(clientName, queueSizeLimit);
#endif
#if defined(__WINDOWS_UWP__)
  if (api == WINDOWS_UWP)
    macapi_ = new MidiInWinUWP(clientName, queueSizeLimit);
#endif
#if defined(__MACOSX_CORE__)
  if (api == MACOSX_CORE)
    macapi_ = new MidiInCore(clientName, queueSizeLimit);
#endif
#if defined(__WEB_MIDI_API__)
  if (api == WEB_MIDI_API)
    macapi_ = new MidiInWeb(clientName, queueSizeLimit);
#endif
#if defined(__AMIDI__)
  if (api == ANDROID_AMIDI)
    macapi_ = new MidiInAndroid(clientName, queueSizeLimit);
#endif
#if defined(__MACMIDI_DUMMY__)
  if (api == MACMIDI_DUMMY)
    macapi_ = new MidiInDummy(clientName, queueSizeLimit);
#endif
}

MACMIDI_DLL_PUBLIC MacMidiIn ::MacMidiIn(MacMidi::Api api,
                                         const std::string &clientName,
                                         unsigned int queueSizeLimit)
    : MacMidi() {
  if (api != UNSPECIFIED) {
    // Attempt to open the specified API.
    openMidiApi(api, clientName, queueSizeLimit);
    if (macapi_)
      return;

    // No compiled support for specified API value.  Issue a warning
    // and continue as if no API was specified.
    std::cerr
        << "\nMacMidiIn: no compiled support for specified API argument!\n\n"
        << std::endl;
  }

  // Iterate through the compiled APIs and return as soon as we find
  // one with at least one port or we reach the end of the list.
  std::vector<MacMidi::Api> apis;
  getCompiledApi(apis);
  for (unsigned int i = 0; i < apis.size(); i++) {
    openMidiApi(apis[i], clientName, queueSizeLimit);
    if (macapi_ && macapi_->getPortCount())
      break;
  }

  if (macapi_)
    return;

  // It should not be possible to get here because the preprocessor
  // definition __MACMIDI_DUMMY__ is automatically defined if no
  // API-specific definitions are passed to the compiler. But just in
  // case something weird happens, we'll throw an error.
  std::string errorText =
      "MacMidiIn: no compiled API support found ... critical error!!";
  throw(MacMidiError(errorText, MacMidiError::UNSPECIFIED));
}

MacMidiIn ::~MacMidiIn() throw() {}

//*********************************************************************//
//  MacMidiOut Definitions
//*********************************************************************//

void MacMidiOut ::openMidiApi(MacMidi::Api api, const std::string &clientName) {
  delete macapi_;
  macapi_ = 0;

#if defined(__UNIX_JACK__)
  if (api == UNIX_JACK)
    macapi_ = new MidiOutJack(clientName);
#endif
#if defined(__LINUX_ALSA__)
  if (api == LINUX_ALSA)
    macapi_ = new MidiOutAlsa(clientName);
#endif
#if defined(__WINDOWS_MM__)
  if (api == WINDOWS_MM)
    macapi_ = new MidiOutWinMM(clientName);
#endif
#if defined(__WINDOWS_UWP__)
  if (api == WINDOWS_UWP)
    macapi_ = new MidiOutWinUWP(clientName);
#endif
#if defined(__MACOSX_CORE__)
  if (api == MACOSX_CORE)
    macapi_ = new MidiOutCore(clientName);
#endif
#if defined(__WEB_MIDI_API__)
  if (api == WEB_MIDI_API)
    macapi_ = new MidiOutWeb(clientName);
#endif
#if defined(__AMIDI__)
  if (api == ANDROID_AMIDI)
    macapi_ = new MidiOutAndroid(clientName);
#endif
#if defined(__MACMIDI_DUMMY__)
  if (api == MACMIDI_DUMMY)
    macapi_ = new MidiOutDummy(clientName);
#endif
}

MACMIDI_DLL_PUBLIC MacMidiOut ::MacMidiOut(MacMidi::Api api,
                                           const std::string &clientName) {
  if (api != UNSPECIFIED) {
    // Attempt to open the specified API.
    openMidiApi(api, clientName);
    if (macapi_)
      return;

    // No compiled support for specified API value.  Issue a warning
    // and continue as if no API was specified.
    std::cerr
        << "\nMacMidiOut: no compiled support for specified API argument!\n\n"
        << std::endl;
  }

  // Iterate through the compiled APIs and return as soon as we find
  // one with at least one port or we reach the end of the list.
  std::vector<MacMidi::Api> apis;
  getCompiledApi(apis);
  for (unsigned int i = 0; i < apis.size(); i++) {
    openMidiApi(apis[i], clientName);
    if (macapi_ && macapi_->getPortCount())
      break;
  }

  if (macapi_)
    return;

  // It should not be possible to get here because the preprocessor
  // definition __MACMIDI_DUMMY__ is automatically defined if no
  // API-specific definitions are passed to the compiler. But just in
  // case something weird happens, we'll thrown an error.
  std::string errorText =
      "MacMidiOut: no compiled API support found ... critical error!!";
  throw(MacMidiError(errorText, MacMidiError::UNSPECIFIED));
}

MacMidiOut ::~MacMidiOut() throw() {}

//*********************************************************************//
//  Common MidiApi Definitions
//*********************************************************************//

MidiApi ::MidiApi(void)
    : apiData_(0), connected_(false), errorCallback_(0),
      firstErrorOccurred_(false), errorCallbackUserData_(0) {}

MidiApi ::~MidiApi(void) {}

void MidiApi ::setErrorCallback(MacMidiErrorCallback errorCallback,
                                void *userData = 0) {
  errorCallback_ = errorCallback;
  errorCallbackUserData_ = userData;
}

void MidiApi ::error(MacMidiError::Type type, std::string errorString) {
  if (errorCallback_) {

    if (firstErrorOccurred_)
      return;

    firstErrorOccurred_ = true;
    const std::string errorMessage = errorString;

    errorCallback_(type, errorMessage, errorCallbackUserData_);
    firstErrorOccurred_ = false;
    return;
  }

  if (type == MacMidiError::WARNING) {
    std::cerr << '\n' << errorString << "\n\n";
  } else if (type == MacMidiError::DEBUG_WARNING) {
#if defined(__MACMIDI_DEBUG__)
    std::cerr << '\n' << errorString << "\n\n";
#endif
  } else {
    std::cerr << '\n' << errorString << "\n\n";
    throw MacMidiError(errorString, type);
  }
}

//*********************************************************************//
//  Common MidiInApi Definitions
//*********************************************************************//

MidiInApi ::MidiInApi(unsigned int queueSizeLimit) : MidiApi() {
  // Allocate the MIDI queue.
  inputData_.queue.ringSize = queueSizeLimit;
  if (inputData_.queue.ringSize > 0)
    inputData_.queue.ring = new MidiMessage[inputData_.queue.ringSize];
}

MidiInApi ::~MidiInApi(void) {
  // Delete the MIDI queue.
  if (inputData_.queue.ringSize > 0)
    delete[] inputData_.queue.ring;
}

void MidiInApi ::setCallback(MacMidiIn::MacMidiCallback callback,
                             void *userData) {
  if (inputData_.usingCallback) {
    errorString_ =
        "MidiInApi::setCallback: a callback function is already set!";
    error(MacMidiError::WARNING, errorString_);
    return;
  }

  if (!callback) {
    errorString_ =
        "MacMidiIn::setCallback: callback function value is invalid!";
    error(MacMidiError::WARNING, errorString_);
    return;
  }

  inputData_.userCallback = callback;
  inputData_.userData = userData;
  inputData_.usingCallback = true;
}

void MidiInApi ::cancelCallback() {
  if (!inputData_.usingCallback) {
    errorString_ = "MacMidiIn::cancelCallback: no callback function was set!";
    error(MacMidiError::WARNING, errorString_);
    return;
  }

  inputData_.userCallback = 0;
  inputData_.userData = 0;
  inputData_.usingCallback = false;
}

void MidiInApi ::ignoreTypes(bool midiSysex, bool midiTime, bool midiSense) {
  inputData_.ignoreFlags = 0;
  if (midiSysex)
    inputData_.ignoreFlags = 0x01;
  if (midiTime)
    inputData_.ignoreFlags |= 0x02;
  if (midiSense)
    inputData_.ignoreFlags |= 0x04;
}

double MidiInApi ::getMessage(std::vector<unsigned char> *message) {
  message->clear();

  if (inputData_.usingCallback) {
    errorString_ =
        "MacMidiIn::getNextMessage: a user callback is currently set "
        "for this port.";
    error(MacMidiError::WARNING, errorString_);
    return 0.0;
  }

  double timeStamp;
  if (!inputData_.queue.pop(message, &timeStamp))
    return 0.0;

  return timeStamp;
}

void MidiInApi ::setBufferSize(unsigned int size, unsigned int count) {
  inputData_.bufferSize = size;
  inputData_.bufferCount = count;
}

unsigned int MidiInApi::MidiQueue::size(unsigned int *__back,
                                        unsigned int *__front) {
  // Access back/front members exactly once and make stack copies for
  // size calculation
  unsigned int _back = back, _front = front, _size;
  if (_back >= _front)
    _size = _back - _front;
  else
    _size = ringSize - _front + _back;

  // Return copies of back/front so no new and unsynchronized accesses
  // to member variables are needed.
  if (__back)
    *__back = _back;
  if (__front)
    *__front = _front;
  return _size;
}

// As long as we haven't reached our queue size limit, push the message.
bool MidiInApi::MidiQueue::push(const MidiInApi::MidiMessage &msg) {
  // Local stack copies of front/back
  unsigned int _back, _front, _size;

  // Get back/front indexes exactly once and calculate current size
  _size = size(&_back, &_front);

  if (_size < ringSize - 1) {
    ring[_back] = msg;
    back = (back + 1) % ringSize;
    return true;
  }

  return false;
}

bool MidiInApi::MidiQueue::pop(std::vector<unsigned char> *msg,
                               double *timeStamp) {
  // Local stack copies of front/back
  unsigned int _back, _front, _size;

  // Get back/front indexes exactly once and calculate current size
  _size = size(&_back, &_front);

  if (_size == 0)
    return false;

  // Copy queued message to the vector pointer argument and then "pop" it.
  msg->assign(ring[_front].bytes.begin(), ring[_front].bytes.end());
  *timeStamp = ring[_front].timeStamp;

  // Update front
  front = (front + 1) % ringSize;
  return true;
}

//*********************************************************************//
//  Common MidiOutApi Definitions
//*********************************************************************//

MidiOutApi ::MidiOutApi(void) : MidiApi() {}

MidiOutApi ::~MidiOutApi(void) {}

// *************************************************** //
//
// OS/API-specific methods.
//
// *************************************************** //

// The CoreMIDI API is based on the use of a callback function for
// MIDI input.  We convert the system specific time stamps to delta
// time values.

// These are not available on iOS.
#if (TARGET_OS_IPHONE == 0)
#include <CoreAudio/HostTime.h>
#include <CoreServices/CoreServices.h>
#endif

// A structure to hold variables related to the CoreMIDI API
// implementation.
struct CoreMidiData {
  MIDIClientRef client;
  MIDIPortRef port;
  MIDIEndpointRef endpoint;
  MIDIEndpointRef destinationId;
  unsigned long long lastTime;
  MIDISysexSendRequest sysexreq;
};

static MIDIClientRef CoreMidiClientSingleton = 0;

void createCoreMidiClientSingleton(const std::string &clientName,
                                   MIDINotifyProc callback) throw() {
  if (CoreMidiClientSingleton == 0) {
    // Set up our client.
    MIDIClientRef client;

    CFStringRef name = CFStringCreateWithCString(NULL, clientName.c_str(),
                                                 kCFStringEncodingASCII);
    OSStatus result = MIDIClientCreate(name, callback, NULL, &client);
    if (result != noErr) {
      std::cerr << "MacMidi: Error creating CoreMIDI client!" << std::endl;
      exit(1);
    }
    CFRelease(name);

    CoreMidiClientSingleton = client;
  }
  return;
}

void MacMidi_setCoreMidiClientSingleton(MIDIClientRef client) {
  CoreMidiClientSingleton = client;
}

void MacMidi_disposeCoreMidiClientSingleton() {
  if (CoreMidiClientSingleton == 0) {
    return;
  }
  MIDIClientDispose(CoreMidiClientSingleton);
  CoreMidiClientSingleton = 0;
}

//*********************************************************************//
//  API: OS-X
//  Class Definitions: MidiInCore
//*********************************************************************//

static void midiInputCallback(const MIDIPacketList *list, void *procRef,
                              void * /*srcRef*/) {
  MidiInApi::MacMidiInData *data =
      static_cast<MidiInApi::MacMidiInData *>(procRef);
  CoreMidiData *apiData = static_cast<CoreMidiData *>(data->apiData);

  unsigned char status;
  unsigned short nBytes, iByte, size;
  unsigned long long time;

  bool &continueSysex = data->continueSysex;
  MidiInApi::MidiMessage &message = data->message;

  const MIDIPacket *packet = &list->packet[0];
  for (unsigned int i = 0; i < list->numPackets; ++i) {

    // My interpretation of the CoreMIDI documentation: all message
    // types, except sysex, are complete within a packet and there may
    // be several of them in a single packet.  Sysex messages can be
    // broken across multiple packets and PacketLists but are bundled
    // alone within each packet (these packets do not contain other
    // message types).  If sysex messages are split across multiple
    // MIDIPacketLists, they must be handled by multiple calls to this
    // function.

    nBytes = packet->length;
    if (nBytes == 0) {
      packet = MIDIPacketNext(packet);
      continue;
    }

    // Calculate time stamp.
    if (data->firstMessage) {
      message.timeStamp = 0.0;
      data->firstMessage = false;
    } else {
      time = packet->timeStamp;
      if (time ==
          0) { // this happens when receiving asynchronous sysex messages
        time = AudioGetCurrentHostTime();
      }
      time -= apiData->lastTime;
      time = AudioConvertHostTimeToNanos(time);
      if (!continueSysex)
        message.timeStamp = time * 0.000000001;
    }

    // Track whether any non-filtered messages were found in this
    // packet for timestamp calculation
    bool foundNonFiltered = false;

    iByte = 0;
    if (continueSysex) {
      // We have a continuing, segmented sysex message.
      if (!(data->ignoreFlags & 0x01)) {
        // If we're not ignoring sysex messages, copy the entire packet.
        for (unsigned int j = 0; j < nBytes; ++j)
          message.bytes.push_back(packet->data[j]);
      }
      continueSysex = packet->data[nBytes - 1] != 0xF7;

      if (!(data->ignoreFlags & 0x01) && !continueSysex) {
        // If not a continuing sysex message, invoke the user callback function
        // or queue the message.
        if (data->usingCallback) {
          MacMidiIn::MacMidiCallback callback =
              (MacMidiIn::MacMidiCallback)data->userCallback;
          callback(message.timeStamp, &message.bytes, data->userData);
        }
        message.bytes.clear();
      }
    } else {
      while (iByte < nBytes) {
        size = 0;
        // We are expecting that the next byte in the packet is a status byte.
        status = packet->data[iByte];
        if (!(status & 0x80))
          break;
        // Determine the number of bytes in the MIDI message.
        if (status < 0xC0)
          size = 3;
        else if (status < 0xE0)
          size = 2;
        else if (status < 0xF0)
          size = 3;
        else if (status == 0xF0) {
          // A MIDI sysex
          if (data->ignoreFlags & 0x01) {
            size = 0;
            iByte = nBytes;
          } else
            size = nBytes - iByte;
          continueSysex = packet->data[nBytes - 1] != 0xF7;
        } else if (status == 0xF1) {
          // A MIDI time code message
          if (data->ignoreFlags & 0x02) {
            size = 0;
            iByte += 2;
          } else
            size = 2;
        } else if (status == 0xF2)
          size = 3;
        else if (status == 0xF3)
          size = 2;
        else if (status == 0xF8 && (data->ignoreFlags & 0x02)) {
          // A MIDI timing tick message and we're ignoring it.
          size = 0;
          iByte += 1;
        } else if (status == 0xFE && (data->ignoreFlags & 0x04)) {
          // A MIDI active sensing message and we're ignoring it.
          size = 0;
          iByte += 1;
        } else
          size = 1;

        // Copy the MIDI data to our vector.
        if (size) {
          foundNonFiltered = true;
          message.bytes.assign(&packet->data[iByte],
                               &packet->data[iByte + size]);
          if (!continueSysex) {
            // If not a continuing sysex message, invoke the user callback
            // function or queue the message.
            if (data->usingCallback) {
              MacMidiIn::MacMidiCallback callback =
                  (MacMidiIn::MacMidiCallback)data->userCallback;
              callback(message.timeStamp, &message.bytes, data->userData);
            }
            message.bytes.clear();
          }
          iByte += size;
        }
      }
    }

    // Save the time of the last non-filtered message
    if (foundNonFiltered) {
      apiData->lastTime = packet->timeStamp;
      if (apiData->lastTime ==
          0) { // this happens when receiving asynchronous sysex messages
        apiData->lastTime = AudioGetCurrentHostTime();
      }
    }

    packet = MIDIPacketNext(packet);
  }
}

MidiInCore ::MidiInCore(const std::string &clientName,
                        unsigned int queueSizeLimit)
    : MidiInApi(queueSizeLimit) {
  MidiInCore::initialize(clientName);
}

MidiInCore ::~MidiInCore(void) {
  // Close a connection if it exists.
  MidiInCore::closePort();

  // Cleanup.
  CoreMidiData *data = static_cast<CoreMidiData *>(apiData_);
  if (data->endpoint)
    MIDIEndpointDispose(data->endpoint);
  delete data;
}

MIDIClientRef
MidiInCore::getCoreMidiClientSingleton(const std::string &clientName) throw() {

  if (CoreMidiClientSingleton == 0) {
    // Set up our client.
    MIDIClientRef client;

    CFStringRef name = CFStringCreateWithCString(NULL, clientName.c_str(),
                                                 kCFStringEncodingASCII);
    OSStatus result = MIDIClientCreate(name, NULL, NULL, &client);
    if (result != noErr) {
      std::ostringstream ost;
      ost << "MidiInCore::initialize: error creating OS-X MIDI client object ("
          << result << ").";
      errorString_ = ost.str();
      error(MacMidiError::DRIVER_ERROR, errorString_);
      return 0;
    }
    CFRelease(name);

    CoreMidiClientSingleton = client;
  }

  return CoreMidiClientSingleton;
}

void MidiInCore ::initialize(const std::string &clientName) {
  // Set up our client.
  MIDIClientRef client = getCoreMidiClientSingleton(clientName);

  // Save our api-specific connection information.
  CoreMidiData *data = (CoreMidiData *)new CoreMidiData;
  data->client = client;
  data->endpoint = 0;
  apiData_ = (void *)data;
  inputData_.apiData = (void *)data;
}

void MidiInCore ::openPort(unsigned int portNumber,
                           const std::string &portName) {
  if (connected_) {
    errorString_ = "MidiInCore::openPort: a valid connection already exists!";
    error(MacMidiError::WARNING, errorString_);
    return;
  }

  CFRunLoopRunInMode(kCFRunLoopDefaultMode, 0, false);
  unsigned int nSrc = MIDIGetNumberOfSources();
  if (nSrc < 1) {
    errorString_ = "MidiInCore::openPort: no MIDI input sources found!";
    error(MacMidiError::NO_DEVICES_FOUND, errorString_);
    return;
  }

  if (portNumber >= nSrc) {
    std::ostringstream ost;
    ost << "MidiInCore::openPort: the 'portNumber' argument (" << portNumber
        << ") is invalid.";
    errorString_ = ost.str();
    error(MacMidiError::INVALID_PARAMETER, errorString_);
    return;
  }

  MIDIPortRef port;
  CoreMidiData *data = static_cast<CoreMidiData *>(apiData_);
  CFStringRef portNameRef =
      CFStringCreateWithCString(NULL, portName.c_str(), kCFStringEncodingASCII);
  OSStatus result = MIDIInputPortCreate(
      data->client, portNameRef, midiInputCallback, (void *)&inputData_, &port);
  CFRelease(portNameRef);

  if (result != noErr) {
    errorString_ = "MidiInCore::openPort: error creating OS-X MIDI input port.";
    error(MacMidiError::DRIVER_ERROR, errorString_);
    return;
  }

  // Get the desired input source identifier.
  MIDIEndpointRef endpoint = MIDIGetSource(portNumber);
  if (endpoint == 0) {
    MIDIPortDispose(port);
    errorString_ =
        "MidiInCore::openPort: error getting MIDI input source reference.";
    error(MacMidiError::DRIVER_ERROR, errorString_);
    return;
  }

  // Make the connection.
  result = MIDIPortConnectSource(port, endpoint, NULL);
  if (result != noErr) {
    MIDIPortDispose(port);
    errorString_ =
        "MidiInCore::openPort: error connecting OS-X MIDI input port.";
    error(MacMidiError::DRIVER_ERROR, errorString_);
    return;
  }

  // Save our api-specific port information.
  data->port = port;

  connected_ = true;
}

void MidiInCore ::openVirtualPort(const std::string &portName) {
  CoreMidiData *data = static_cast<CoreMidiData *>(apiData_);

  // Create a virtual MIDI input destination.
  MIDIEndpointRef endpoint;
  CFStringRef portNameRef =
      CFStringCreateWithCString(NULL, portName.c_str(), kCFStringEncodingASCII);
  OSStatus result =
      MIDIDestinationCreate(data->client, portNameRef, midiInputCallback,
                            (void *)&inputData_, &endpoint);
  CFRelease(portNameRef);

  if (result != noErr) {
    errorString_ = "MidiInCore::openVirtualPort: error creating virtual OS-X "
                   "MIDI destination.";
    error(MacMidiError::DRIVER_ERROR, errorString_);
    return;
  }

  // Save our api-specific connection information.
  data->endpoint = endpoint;
}

void MidiInCore ::closePort(void) {
  CoreMidiData *data = static_cast<CoreMidiData *>(apiData_);

  if (data->endpoint) {
    MIDIEndpointDispose(data->endpoint);
    data->endpoint = 0;
  }

  if (data->port) {
    MIDIPortDispose(data->port);
    data->port = 0;
  }

  connected_ = false;
}

void MidiInCore ::setClientName(const std::string &) {

  errorString_ = "MidiInCore::setClientName: this function is not implemented "
                 "for the MACOSX_CORE API!";
  error(MacMidiError::WARNING, errorString_);
}

void MidiInCore ::setPortName(const std::string &) {

  errorString_ = "MidiInCore::setPortName: this function is not implemented "
                 "for the MACOSX_CORE API!";
  error(MacMidiError::WARNING, errorString_);
}

unsigned int MidiInCore ::getPortCount() {
  CFRunLoopRunInMode(kCFRunLoopDefaultMode, 0, false);
  return MIDIGetNumberOfSources();
}

// This function was submitted by Douglas Casey Tucker and apparently
// derived largely from PortMidi.
CFStringRef CreateEndpointName(MIDIEndpointRef endpoint, bool isExternal) {
  CFMutableStringRef result = CFStringCreateMutable(NULL, 0);
  CFStringRef str;

  // Begin with the endpoint's name.
  str = NULL;
  MIDIObjectGetStringProperty(endpoint, kMIDIPropertyName, &str);
  if (str != NULL) {
    CFStringAppend(result, str);
  }

  // some MIDI devices have a leading space in endpoint name. trim
  CFStringTrim(result, CFSTR(" "));

  MIDIEntityRef entity = 0;
  MIDIEndpointGetEntity(endpoint, &entity);
  if (entity == 0)
    // probably virtual
    return result;

  if (CFStringGetLength(result) == 0) {
    // endpoint name has zero length -- try the entity
    str = NULL;
    MIDIObjectGetStringProperty(entity, kMIDIPropertyName, &str);
    if (str != NULL) {
      CFStringAppend(result, str);
    }
  }
  // now consider the device's name
  MIDIDeviceRef device = 0;
  MIDIEntityGetDevice(entity, &device);
  if (device == 0)
    return result;

  str = NULL;
  MIDIObjectGetStringProperty(device, kMIDIPropertyName, &str);
  if (CFStringGetLength(result) == 0) {
    CFRelease(result);
    CFRetain(str);
    return str;
  }
  if (str != NULL) {
    // if an external device has only one entity, throw away
    // the endpoint name and just use the device name
    if (isExternal && MIDIDeviceGetNumberOfEntities(device) < 2) {
      CFRelease(result);
      CFRetain(str);
      return str;
    } else {
      if (CFStringGetLength(str) == 0) {
        return result;
      }
      // does the entity name already start with the device name?
      // (some drivers do this though they shouldn't)
      // if so, do not prepend
      if (CFStringCompareWithOptions(result, /* endpoint name */
                                     str /* device name */,
                                     CFRangeMake(0, CFStringGetLength(str)),
                                     0) != kCFCompareEqualTo) {
        // prepend the device name to the entity name
        if (CFStringGetLength(result) > 0)
          CFStringInsert(result, 0, CFSTR(" "));

        CFStringInsert(result, 0, str);
      }
    }
  }
  return result;
}

// This function was submitted by Douglas Casey Tucker and apparently
// derived largely from PortMidi.
static CFStringRef CreateConnectedEndpointName(MIDIEndpointRef endpoint) {
  CFMutableStringRef result = CFStringCreateMutable(NULL, 0);
  CFStringRef str;
  OSStatus err;
  int i;

  // Does the endpoint have connections?
  CFDataRef connections = NULL;
  int nConnected = 0;
  bool anyStrings = false;
  err = MIDIObjectGetDataProperty(endpoint, kMIDIPropertyConnectionUniqueID,
                                  &connections);
  if (connections != NULL) {
    // It has connections, follow them
    // Concatenate the names of all connected devices
    nConnected = CFDataGetLength(connections) / sizeof(MIDIUniqueID);
    if (nConnected) {
      const SInt32 *pid = (const SInt32 *)(CFDataGetBytePtr(connections));
      for (i = 0; i < nConnected; ++i, ++pid) {
        MIDIUniqueID id = EndianS32_BtoN(*pid);
        MIDIObjectRef connObject;
        MIDIObjectType connObjectType;
        err = MIDIObjectFindByUniqueID(id, &connObject, &connObjectType);
        if (err == noErr) {
          if (connObjectType == kMIDIObjectType_ExternalSource ||
              connObjectType == kMIDIObjectType_ExternalDestination) {
            // Connected to an external device's endpoint (10.3 and later).
            str = CreateEndpointName((MIDIEndpointRef)(connObject), true);
          } else {
            // Connected to an external device (10.2) (or something else, catch-
            str = NULL;
            MIDIObjectGetStringProperty(connObject, kMIDIPropertyName, &str);
            if (str)
              CFRetain(str);
          }
          if (str != NULL) {
            if (anyStrings)
              CFStringAppend(result, CFSTR(", "));
            else
              anyStrings = true;
            CFStringAppend(result, str);
            CFRelease(str);
          }
        }
      }
    }
    CFRelease(connections);
  }
  if (anyStrings)
    return result;

  CFRelease(result);

  // Here, either the endpoint had no connections, or we failed to obtain names
  return CreateEndpointName(endpoint, false);
}

std::string MidiInCore ::getPortName(unsigned int portNumber) {
  CFStringRef nameRef;
  MIDIEndpointRef portRef;
  char name[128];

  std::string stringName;
  CFRunLoopRunInMode(kCFRunLoopDefaultMode, 0, false);
  if (portNumber >= MIDIGetNumberOfSources()) {
    std::ostringstream ost;
    ost << "MidiInCore::getPortName: the 'portNumber' argument (" << portNumber
        << ") is invalid.";
    errorString_ = ost.str();
    error(MacMidiError::WARNING, errorString_);
    return stringName;
  }

  portRef = MIDIGetSource(portNumber);
  nameRef = CreateConnectedEndpointName(portRef);
  CFStringGetCString(nameRef, name, sizeof(name), kCFStringEncodingUTF8);
  CFRelease(nameRef);

  return stringName = name;
}

//*********************************************************************//
//  API: OS-X
//  Class Definitions: MidiOutCore
//*********************************************************************//

MidiOutCore ::MidiOutCore(const std::string &clientName) : MidiOutApi() {
  MidiOutCore::initialize(clientName);
}

MidiOutCore ::~MidiOutCore(void) {
  // Close a connection if it exists.
  MidiOutCore::closePort();

  // Cleanup.
  CoreMidiData *data = static_cast<CoreMidiData *>(apiData_);
  if (data->endpoint)
    MIDIEndpointDispose(data->endpoint);
  delete data;
}

MIDIClientRef
MidiOutCore::getCoreMidiClientSingleton(const std::string &clientName) throw() {

  if (CoreMidiClientSingleton == 0) {
    // Set up our client.
    MIDIClientRef client;

    CFStringRef name = CFStringCreateWithCString(NULL, clientName.c_str(),
                                                 kCFStringEncodingASCII);
    OSStatus result = MIDIClientCreate(name, NULL, NULL, &client);
    if (result != noErr) {
      std::ostringstream ost;
      ost << "MidiInCore::initialize: error creating OS-X MIDI client object ("
          << result << ").";
      errorString_ = ost.str();
      error(MacMidiError::DRIVER_ERROR, errorString_);
      return 0;
    }
    CFRelease(name);

    CoreMidiClientSingleton = client;
  }

  return CoreMidiClientSingleton;
}

void MidiOutCore ::initialize(const std::string &clientName) {
  // Set up our client.
  MIDIClientRef client = getCoreMidiClientSingleton(clientName);

  // Save our api-specific connection information.
  CoreMidiData *data = (CoreMidiData *)new CoreMidiData;
  data->client = client;
  data->endpoint = 0;
  apiData_ = (void *)data;
}

unsigned int MidiOutCore ::getPortCount() {
  CFRunLoopRunInMode(kCFRunLoopDefaultMode, 0, false);
  return MIDIGetNumberOfDestinations();
}

std::string MidiOutCore ::getPortName(unsigned int portNumber) {
  CFStringRef nameRef;
  MIDIEndpointRef portRef;
  char name[128];

  std::string stringName;
  CFRunLoopRunInMode(kCFRunLoopDefaultMode, 0, false);
  if (portNumber >= MIDIGetNumberOfDestinations()) {
    std::ostringstream ost;
    ost << "MidiOutCore::getPortName: the 'portNumber' argument (" << portNumber
        << ") is invalid.";
    errorString_ = ost.str();
    error(MacMidiError::WARNING, errorString_);
    return stringName;
  }

  portRef = MIDIGetDestination(portNumber);
  nameRef = CreateConnectedEndpointName(portRef);
  CFStringGetCString(nameRef, name, sizeof(name), kCFStringEncodingUTF8);
  CFRelease(nameRef);

  return stringName = name;
}

void MidiOutCore ::openPort(unsigned int portNumber,
                            const std::string &portName) {
  if (connected_) {
    errorString_ = "MidiOutCore::openPort: a valid connection already exists!";
    error(MacMidiError::WARNING, errorString_);
    return;
  }

  CFRunLoopRunInMode(kCFRunLoopDefaultMode, 0, false);
  unsigned int nDest = MIDIGetNumberOfDestinations();
  if (nDest < 1) {
    errorString_ = "MidiOutCore::openPort: no MIDI output destinations found!";
    error(MacMidiError::NO_DEVICES_FOUND, errorString_);
    return;
  }

  if (portNumber >= nDest) {
    std::ostringstream ost;
    ost << "MidiOutCore::openPort: the 'portNumber' argument (" << portNumber
        << ") is invalid.";
    errorString_ = ost.str();
    error(MacMidiError::INVALID_PARAMETER, errorString_);
    return;
  }

  MIDIPortRef port;
  CoreMidiData *data = static_cast<CoreMidiData *>(apiData_);
  CFStringRef portNameRef =
      CFStringCreateWithCString(NULL, portName.c_str(), kCFStringEncodingASCII);
  OSStatus result = MIDIOutputPortCreate(data->client, portNameRef, &port);
  CFRelease(portNameRef);
  if (result != noErr) {
    errorString_ =
        "MidiOutCore::openPort: error creating OS-X MIDI output port.";
    error(MacMidiError::DRIVER_ERROR, errorString_);
    return;
  }

  // Get the desired output port identifier.
  MIDIEndpointRef destination = MIDIGetDestination(portNumber);
  if (destination == 0) {
    MIDIPortDispose(port);
    errorString_ = "MidiOutCore::openPort: error getting MIDI output "
                   "destination reference.";
    error(MacMidiError::DRIVER_ERROR, errorString_);
    return;
  }

  // Save our api-specific connection information.
  data->port = port;
  data->destinationId = destination;
  connected_ = true;
}

void MidiOutCore ::closePort(void) {
  CoreMidiData *data = static_cast<CoreMidiData *>(apiData_);

  if (data->endpoint) {
    MIDIEndpointDispose(data->endpoint);
    data->endpoint = 0;
  }

  if (data->port) {
    MIDIPortDispose(data->port);
    data->port = 0;
  }

  connected_ = false;
}

void MidiOutCore ::setClientName(const std::string &) {

  errorString_ = "MidiOutCore::setClientName: this function is not implemented "
                 "for the MACOSX_CORE API!";
  error(MacMidiError::WARNING, errorString_);
}

void MidiOutCore ::setPortName(const std::string &) {

  errorString_ = "MidiOutCore::setPortName: this function is not implemented "
                 "for the MACOSX_CORE API!";
  error(MacMidiError::WARNING, errorString_);
}

void MidiOutCore ::openVirtualPort(const std::string &portName) {
  CoreMidiData *data = static_cast<CoreMidiData *>(apiData_);

  if (data->endpoint) {
    errorString_ =
        "MidiOutCore::openVirtualPort: a virtual output port already exists!";
    error(MacMidiError::WARNING, errorString_);
    return;
  }

  // Create a virtual MIDI output source.
  MIDIEndpointRef endpoint;
  CFStringRef portNameRef =
      CFStringCreateWithCString(NULL, portName.c_str(), kCFStringEncodingASCII);
  OSStatus result = MIDISourceCreate(data->client, portNameRef, &endpoint);
  CFRelease(portNameRef);

  if (result != noErr) {
    errorString_ =
        "MidiOutCore::initialize: error creating OS-X virtual MIDI source.";
    error(MacMidiError::DRIVER_ERROR, errorString_);
    return;
  }

  // Save our api-specific connection information.
  data->endpoint = endpoint;
}

void MidiOutCore ::sendMessage(const unsigned char *message, size_t size) {
  // We use the MIDISendSysex() function to asynchronously send sysex
  // messages.  Otherwise, we use a single CoreMidi MIDIPacket.
  unsigned int nBytes = static_cast<unsigned int>(size);
  if (nBytes == 0) {
    errorString_ = "MidiOutCore::sendMessage: no data in message argument!";
    error(MacMidiError::WARNING, errorString_);
    return;
  }

  if (message[0] != 0xF0 && nBytes > 3) {
    errorString_ = "MidiOutCore::sendMessage: message format problem ... not "
                   "sysex but > 3 bytes?";
    error(MacMidiError::WARNING, errorString_);
    return;
  }

  MIDITimeStamp timeStamp = AudioGetCurrentHostTime();
  CoreMidiData *data = static_cast<CoreMidiData *>(apiData_);
  OSStatus result;

  ByteCount bufsize = nBytes > 65535 ? 65535 : nBytes;
  Byte buffer[bufsize + 16]; // pad for other struct members
  ByteCount listSize = sizeof(buffer);
  MIDIPacketList *packetList = (MIDIPacketList *)buffer;

  ByteCount remainingBytes = nBytes;
  while (remainingBytes) {
    MIDIPacket *packet = MIDIPacketListInit(packetList);
    // A MIDIPacketList can only contain a maximum of 64K of data, so if our
    // message is longer, break it up into chunks of 64K or less and send out as
    // a MIDIPacketList with only one MIDIPacket. Here, we reuse the memory
    // allocated above on the stack for all.
    ByteCount bytesForPacket = remainingBytes > 65535 ? 65535 : remainingBytes;
    const Byte *dataStartPtr = (const Byte *)&message[nBytes - remainingBytes];
    packet = MIDIPacketListAdd(packetList, listSize, packet, timeStamp,
                               bytesForPacket, dataStartPtr);
    remainingBytes -= bytesForPacket;

    if (!packet) {
      errorString_ = "MidiOutCore::sendMessage: could not allocate packet list";
      error(MacMidiError::DRIVER_ERROR, errorString_);
      return;
    }

    // Send to any destinations that may have connected to us.
    if (data->endpoint) {
      result = MIDIReceived(data->endpoint, packetList);
      if (result != noErr) {
        errorString_ = "MidiOutCore::sendMessage: error sending MIDI to "
                       "virtual destinations.";
        error(MacMidiError::WARNING, errorString_);
      }
    }

    // And send to an explicit destination port if we're connected.
    if (connected_) {
      result = MIDISend(data->port, data->destinationId, packetList);
      if (result != noErr) {
        errorString_ =
            "MidiOutCore::sendMessage: error sending MIDI message to port.";
        error(MacMidiError::WARNING, errorString_);
      }
    }
  }
}
