let previousHash = 0;
let ran = false;

const encoder = new TextEncoder();

async function getHash(data) {
  const bytes = encoder.encode(data)
  const hash = await window.crypto.subtle.digest("SHA-256", bytes)
  const hashArray = Array.from(new Uint8Array(hash));
  return hashHex = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
}

setInterval(() => {
  fetch(window.location.pathname)
    .then(response => response.text())
    .then(getHash)
    .then(hash => {
      if (!ran) {
        ran = true;
        previousHash = hash
        return;
      }
      if (hash != previousHash) {
        console.log("Reloading page")
        location.reload();
      }
    })
}, 1000);
