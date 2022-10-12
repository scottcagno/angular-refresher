export class Book {
  title !:string;
  author !:string;
  price !:number;

  constructor(title ?:string, author ?:string, price ?:number) {
    if (title) { this.title = title }
    if (author) { this.author = author }
    if (price) { this.price = price }
  }
}
