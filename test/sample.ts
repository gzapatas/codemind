export class Person {
  name: string;
  age: number;

  constructor(name: string, age: number) {
    this.name = name;
    this.age = age;
  }

  greet(): string {
    return `hello, ${this.name}`;
  }
}

export function add(a: number, b: number): number {
  return a + b;
}
