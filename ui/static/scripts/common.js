class Queue {
    constructor() {
        this.queue = [];
        return this;
    }

    dequeue() {
        let a = this.queue.shift();
        return a;
    }
    enqueue(item) {
        this.queue.push(item);
    }

    // isEmpty returns true if the queue is empty or false if not
    isEmpty() {
        return this.queue.length == 0;
    }
    from(array) {
        if (array != undefined) {
            array.forEach((val) => {
                this.enqueue(val);
            });
            return this;
        }
    }

    get() {
        return this.queue;
    }
}

class connectionPool {
    constructor(ports) {
        this.availableConns = new Queue().from(ports);
        this.connsInUse = new Map();
    }

    getConnection() {
        let available = this.availableConns.dequeue();
        this.connsInUse.set(available, true);
        if (available == undefined) {
            return null;
        }
        return available;
    }

    hasAvailable() {
        return !this.availableConns.isEmpty();
    }

    done(port) {
        this.connsInUse.set(port, false);
        this.availableConns.enqueue(port);
    }
}

function uiUpdate(update) {
    switch (update.uiPortion) {
        case "title":
            let updatedContent = document.createElement("p");
            updatedContent.innerHTML = update.content;
            props.get("uiTitle").innerHTML = "";
            props.get("uiTitle").appendChild(updatedContent);
            break;
        case "progress":
            props.get("progressArea").appendChild(update.content);
            break;
        case "notify-ok":
            // insert notification to document
            break;
        case "notify-success":
            // insert notification to document
            break;
        case "notify-err":
            // insert notification to document
            break;
        default:
            break;
    }
}

function handleStatusChange(statusObj) {
    /* 
      const uiUpdateSpec = {
          uiPortion: "title" | "progress" | "notify-ok" | "noify-success" | "notify-err",
          content: null
      } 
    */
    switch (statusObj.status) {
        case "connected":
            //update UI
            uiUpdate({
                uiPortion: "title",
                content: `Connected to ${statusObj.hostname}`,
            });

            // setup connection pool
            pool = new connectionPool(statusObj.connectionPool);
            props.set("connectionPool", pool);

            // setup connedNodeIP
            props.set("connectedNodeIp", statusObj.connectedIP);

            // append file input
            let fileInput = document.createElement("input");
            fileInput.setAttribute("type", "file");
            fileInput.setAttribute("multiple", true);
            fileInput.onchange = onFilesInputHandler;
            props.get("uiTitle").appendChild(fileInput);

            break;
    }
}

class progressBar {
    constructor() {
        const barContainer = document.createElement("div");
        barContainer.setAttribute("aria-valuenow", "0%");
        barContainer.setAttribute("role", "progressbar");
        const bar = document.createElement("div");
        bar.classList.add("progress-bar");

        const progressValue = document.createElement("p");
        progressValue.classList.add("count");
        barContainer.appendChild(progressValue);
        barContainer.appendChild(bar);

        this.counter = progressValue;
        this.bar = bar;
        this.progressBar = barContainer;
    }

    update(newPercentValue) {
        this.progressBar.setAttribute("aria-valuenow", newPercentValue);
        (function draw() {
            if (i <= 100) {
                requestAnimationFrame(draw);
                bar.style.width = newPercentValue + "%";
                counter.innerHTML = Math.round(newPercentValue) + "%";
            } else {
                bar.className += " done";
            }
        })();
    }
}

class fileHandler {
    static fileQueue = new Queue();

    constructor(file) {
        this.file = file;
        this.progressBar = new progressBar();
    }

    async openConnection(port, chunkSize, chunkCount, filename) {
        // format == ws://ip:port/{chunkSize}-{chunkCount}-{filename}
        if (!props.has("connectedNodeIp")) return null;
        const connectionString = `ws://${props.get(
            "connectedNodeIp"
        )}:${port}/${chunkSize}-${chunkCount}-${filename}`;

        let conn = await new WebSocket(connectionString);
        console.log(conn);
        if (!!conn) return conn;
        else return null;
    }

    send() {
        // request port
        if (!props.has("connectionPool")) {
            /* 
			const uiUpdateSpec = {
				uiPortion: "title" | "progress" | "notify-ok" | "noify-success" | "notify-err",
				content: null
			} 
			*/
            uiUpdate({
                uiPortion: "notify-err",
                content:
                    "Unable to Obtain establish connection: Connection pool not available",
            });
            return;
        }
        let port = props.get("connectionPool").getConnection();
        if (port == null) {
            console.log("adding file to wait queue");
            fileHandler.fileQueue.enqueue(this);
            return;
        }

        const fr = new FileReader();
        fr.onload = async (e) => {
            const content = e.target.result;
            // chunk up file
            const CHUNK_SIZE = Math.trunc(content.byteLength * 0.1);
            console.log("file size ", content.byteLength);
            console.log("chunk size ", CHUNK_SIZE);
            const totalChunks = Math.ceil(content.byteLength / CHUNK_SIZE);
            let writtenChunks = 0;

            const conn = await this.openConnection(
                port,
                CHUNK_SIZE,
                totalChunks,
                this.file.name
            );

            conn.onopen = () => {
                /* 
                    insert progress bar
                */
                // uiUpdate({ uiPortion: "progress", content: this.progressBar.progressBar });
                console.log("Connection open, readystate: ", conn.readyState);
                //send fileChunks
                for (let chunk = 0; chunk < totalChunks; chunk++) {
                    let CHUNK = content.slice(
                        chunk * CHUNK_SIZE,
                        (chunk + 1) * CHUNK_SIZE
                    );
                    conn.send(CHUNK);

                    // update progressbar
                    writtenChunks += CHUNK.byteLength;
                    // this.progressBar.update(
                    //     `${(100 * writtenChunks) / content.byteLength}%`
                    // );
                }
                // after sending close connection and return port to pool
                conn.close();
            };

            conn.onclose = function () {
                props.get("connectionPool").done(port);
                console.info(
                    writtenChunks,
                    " chunks / ",
                    totalChunks,
                    " total chumks written"
                );

                // call done event on window so that filehandler can empty wait queue
                window.dispatchEvent(fileSentEvent);
            };
        };
        fr.readAsArrayBuffer(this.file);
    }
}

const fileSentEvent = new Event("fileSent");
const props = new Map();

props.set("uiTitle", document.getElementById("uiTitle"));
props.set("progressArea", document.getElementById("progressArea"));

window.addEventListener("fileSent", (e) => {
    // clear out filehandler.filequeue
    while (!fileHandler.fileQueue.isEmpty()) {
        if (pool.hasAvailable()) {
            fh = fileHandler.fileQueue.dequeue();
            fh.send();
        } else {
            break;
        }
    }
});

function onFilesInputHandler(e) {
    // get files
    files = e.target.files;
    // for ech file, create a file handler and let it handle
    Array.from(files).forEach((file) => {
        new fileHandler(file).send();
    });
}
