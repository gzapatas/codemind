class Person {
  constructor(name, age) {
    this.name = name;
    this.age = age;
  }

  greet() {
    return `hello, ${this.name}`;
  }
}

function add(a, b) {
  return a + b;
}

module.exports = { Person, add };
