import React from 'react';

export class Person extends React.Component {
  constructor(props) {
    super(props);
  }

  greet() {
    return `hello, ${this.props.name}`;
  }

  render() {
    return <div>{this.greet()}</div>;
  }
}

export function add2(a, b) {
  return a + b;
}

export default Person;
