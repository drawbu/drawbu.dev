let previousHash = undefined;
let previousUrl = window.location.pathname
let ran = false;

const encoder = new TextEncoder();

async function getHash(data) {
  const bytes = encoder.encode(data)
  const hash = await window.crypto.subtle.digest("SHA-256", bytes)
  const hashArray = Array.from(new Uint8Array(hash));
  return hashHex = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
}

setInterval(() => {
  if (window.location.pathname !== previousUrl) {
    previousHash = undefined;
  }

  fetch(window.location.pathname)
    .then(response => response.text())
    .then(getHash)
    .then(hash => {
      if (previousHash === undefined) {
        previousHash = hash
        return;
      }
      if (hash != previousHash) {
        console.log("Reloading page")
        location.reload();
      }
    })
}, 1000);
