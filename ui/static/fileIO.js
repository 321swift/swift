class Queue {
	constructor() {
		this.queue = [];
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

	get() {
		return this.queue;
	}
}
let fileQueue = new Queue();
let fileInput = document.getElementById("fileInput");
let enqueueEvent = new Event("enqueue");

fileInput.addEventListener("change", (event) => {
	let filePromises = [];
	Array.from(event.target.files).forEach((file) => {
		filePromises.push(readFile(file));
	});
	Promise.all(filePromises).then(() => {
		console.log(fileQueue.get());
		window.dispatchEvent(enqueueEvent);
	});
});
window.addEventListener("enqueue", (event) => {
	while (!fileQueue.isEmpty()) {
		if (props.fileSocket != undefined) {
			props.fileSocket.send(JSON.stringify(fileQueue.dequeue()));
		}
	}
});

function readFile(file) {
	return new Promise((resolve, reject) => {
		let fr = new FileReader();
		fr.onload = (event) => {
			fileQueue.enqueue({
				Filename: file.name,
				Data: event.target.result,
			});
			resolve();
		};
		fr.readAsArrayBuffer(file);
	});
}
