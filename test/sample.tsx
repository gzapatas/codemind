import React from 'react';

export class Person extends React.Component<{ name: string; age: number }> {
  constructor(props: { name: string; age: number }) {
    super(props);
  }

  greet(): string {
    return `hello, ${this.props.name}`;
  }

  render() {
    return <div>{this.greet()}</div>;
  }
}

export function add(a: number, b: number): number {
  return a + b;
}

export default Person;
