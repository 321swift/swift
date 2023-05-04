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

	done(port) {
		this.connsInUse.set(port, false);
		this.availableConns.enqueue(port);
	}
}

const fileSentEvent = new Event("fileSent");
let pool = new connectionPool([]);
let props = new Map();

class fileHandler {
	static fileQueue = new Queue();

	constructor(file) {
		this.file = file;
	}

	openConnection(port, chunkSize, chunkCount, filename) {
		// format == ws://ip:port/{chunkSize}-{chunkCount}-{filename}
		connectionString = `ws://${props.get(
			connectedNodeIP
		)}:${port}/${chunkSize}-${chunkCount}-${filename}`;

		let conn = new WebSocket(connectionString);
		return conn;
	}

	send() {
		// request port
		let port = pool.getConnection();
		if (port == null) {
			console.log("adding file to wait queue");
			fileHandler.fileQueue.enqueue(this);
			return;
		}

		fr = new FileReader();
		this.fr.onload = (e) => {
			const content = e.target.result;
			// chunk up file
			const CHUNK_SIZE = 1000;
			const totalChunks = e.target.result.byteLength / CHUNK_SIZE;

			conn = this.openConnection(port, CHUNK_SIZE, chunkCount, this.file.name);

			// loop over each chunk
			for (let chunk = 0; chunk < totalChunks + 1; chunk++) {
				let CHUNK = content.slice(chunk * CHUNK_SIZE, (chunk + 1) * CHUNK_SIZE);
				conn.send(CHUNK);
				//update progressbar
			}
			// after sending close connection and return port to pool
			conn.close();
			pool.done(port);

			// call done event on window so that filehandler can empty wait queue
			window.dispatchEvent(fileSentEvent);
		};
		fr.readAsArrayBuffer(this.file);
	}
}
