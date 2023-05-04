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
