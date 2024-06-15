// @ts-check

/** *
 * @typedef {{
 *  boosted: unknown
 *  elt: HTMLElement
 *  etx: unknown
 *  pathInfo: PathInfo
 *  requestConfig: RequestConfig
 *  select: unknown
 *  target: HTMLElement
 *  xhr: XMLHttpRequest
 * }} EventDetail
 * @typedef {{
 *  requestPath: string
 *  finalRequestPath: string
 *  anchor: unknown
 * }} PathInfo
 * @typedef {{
 *  boosted: unknown
 *  elt: HTMLElement
 *  errors: unknown[]
 *  headers: Record<string, string|null>
 *  parameters: unknown
 *  path: string
 *  target: HTMLElement
 *  timeout: number
 *  triggeringEvent: Event
 *  unfilteredParameters: unknown
 *  useUrlParams: boolean
 *  verb: string
 *  withCredentials: boolean
 * }} RequestConfig
 *
 */

/** @type {any} */
// @ts-ignore
htmx;

document.body.addEventListener("htmx:beforeRequest", function (evt) {
  /** @type {EventDetail} */
  // @ts-ignore
  const detail = evt.detail;
  if (detail.elt.getAttribute("hx-stream") !== "true") {
    return;
  }
  evt.preventDefault();

  console.log(detail);
  stream(detail);
});

/**
 *
 * @param {EventDetail} requestConfig
 */
async function stream({ elt, requestConfig, target }) {
  const abort = new AbortController();
  function onAbort() {
    abort.abort();
  }

  const observer = new MutationObserver((mutations) => {
    for (const mut of mutations) {
      for (const newElt of mut.addedNodes) {
        // @ts-ignore
        htmx.process(newElt);
      }
    }
  });

  try {
    elt.addEventListener("htmx:abort", onAbort);

    observer.observe(target, {
      attributes: true,
      childList: true,
      subtree: true,
    });

    const resp = await fetch(requestConfig.path, {
      method: requestConfig.verb,
      headers: processHeaders(requestConfig.headers),
      // credentials: requestConfig.withCredentials ? 'include' : 'omit',
      signal: abort.signal,
    });
    if (!resp.body) {
      throw new Error("no body on scan");
    }

    if (false) {
    }

    target.innerHTML = "";

    const reader = resp.body.getReader();

    const decoder = new TextDecoder();
    for await (const chunk of read(reader)) {
      target.insertAdjacentHTML("beforeend", decoder.decode(chunk));
    }
  } catch (e) {
    if (e instanceof DOMException) {
      if (e.name === "AbortError") {
        return;
      }
    }
    throw e;
  } finally {
    elt.removeEventListener("htmx:abort", onAbort);
    observer.disconnect();
  }
}

/**
 *
 * @param {ReadableStreamDefaultReader<Uint8Array>} reader
 * @returns {AsyncGenerator<Uint8Array, void>}
 */
async function* read(reader) {
  while (true) {
    const { value: chunk, done } = await reader.read();
    if (done) {
      return;
    }
    yield chunk;
  }
}

/**
 *
 * @param {Record<string,string|null>} headers
 * @returns {Record<string,string>}
 */
function processHeaders(headers) {
  // @ts-ignore
  return Object.fromEntries(
    Object.entries(headers).filter((header) => header[1] !== null)
  );
}
