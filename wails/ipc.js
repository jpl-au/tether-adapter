// Gives tether its socket under mode.IPC: the app's Wails stream,
// wrapped so it behaves as the WebSocket tether expects. Nothing leaves
// the app: Wails carries every frame between the page and Go.
import { Stream } from "/wails/runtime.js";

window.Tether = window.Tether || {};

window.Tether.ipc = function (url) {
  var stream = Stream("tether");
  var text = new TextDecoder();
  var socket = {
    binaryType: "arraybuffer",
    onopen: null,
    onmessage: null,
    onclose: null,
    onerror: null,
    get readyState() { return stream.readyState; },
    send: function (data) { stream.send(data); },
    close: function (code, reason) { stream.close(code, reason); },
  };
  stream.addEventListener("open", function (e) {
    // The first frame says what a WebSocket's request would: the URL
    // with the page's ticket, and the User-Agent the ticket is bound
    // to. It goes before tether's own open handler sends anything.
    stream.send(JSON.stringify({ url: url, userAgent: navigator.userAgent }));
    if (socket.onopen) socket.onopen(e);
  });
  stream.addEventListener("message", function (e) {
    // Wails delivers every frame as bytes. Tether reads text, and a
    // binary wire format reaches it as base64 text, as it does over SSE.
    if (socket.onmessage) socket.onmessage({ data: text.decode(e.data) });
  });
  stream.addEventListener("close", function (e) {
    if (socket.onclose) socket.onclose(e);
  });
  stream.addEventListener("error", function (e) {
    if (socket.onerror) socket.onerror(e);
  });
  return socket;
};
