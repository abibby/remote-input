// @ts-check

async function scan() {
    console.log('scan');
    const resp = await fetch("/bluetooth/scan")
    if (!resp.body){
        throw new Error("no body on scan")
    }
    const reader = resp.body.getReader()
  
    while (true) {
        const { value, done } = await reader.read();
        if (done) {
            return;
        }
        console.log(value);
    }
}

document.getElementById('scan')?.addEventListener('click', e => {
    scan()
})