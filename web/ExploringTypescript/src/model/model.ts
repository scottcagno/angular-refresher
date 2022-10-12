
// method signature, by default is public
//
// <private> methodName(paramName :paramType) :returnType
//
// (note, "?:" for optional parameters)


export enum SubjectArea {
  ART,
  HISTORY,
  SCIENCE,
  LITERATURE,
}

// Book
export class Book {

  readonly id: number;
  title !:string;
  author !:string
  price !:number;

  constructor(id :number) {
    this.id = id;
  }

  stringify() {
    return this.toString()
  }

  toString() :string {
    return `id=${this.id}, title=${this.title}, author=${this.author}, price=${this.price} (no tax)`;
  }

  priceWithTax(taxRate :number) :number {
    return this.price * (1+taxRate)
  }
}

// Video
export class Video {

  private title!: string;
  private author!: string
  private price!: number;

  constructor(title?:string, author?:string, price?:number) {
    if (title) { this.title = title }
    if (author) { this.author = author }
    if (price) { this.price = price }
  }
}
