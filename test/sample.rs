struct Person {
    name: String,
    age: u32,
}

impl Person {
    fn greet(&self) -> String {
        format!("hello, {}", self.name)
    }
}

fn add(a: i32, b: i32) -> i32 {
    a + b
}

fn main() {
    let p = Person { name: "Bob".to_string(), age: 28 };
    println!("{}", p.greet());
    println!("{}", add(1, 2));
}
// crea multiples funciones de test para cada chunker registrado, y cada una debe probar que el chunker devuelve los chunks esperados para un input dado.