export namespace main {
	
	export class Person {
	    name: string;
	    age: number;
	
	    static createFrom(source: any = {}) {
	        return new Person(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.age = source["age"];
	    }
	}

}

