import { Component } from '@angular/core';
import {Book, SubjectArea, Video} from "../model/model";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'ExploringTypescript';

  myNumber!: number;

  readonly myReadOnlyNumber: number = 3;

  constructor() {

    console.log(SubjectArea.HISTORY) // prints enum value
    console.log(SubjectArea[2]) // prints the enum value lable

    //super(props);
    this.exploringArrays();

    // "class" in JavaScript
    let myCustomer = {firstName: 'Matt', age: 21}
    console.log(myCustomer);
    console.log(typeof myCustomer);

    // import an external class
    let myBook = new Book(1);
    let myVideo : Video;

    myBook.author = 'Jon Doe'
    myBook.title = 'A fantastic read'
    console.log(myBook);
    myBook.price = 100;
    console.log('to buy this book it will cost ' + myBook.priceWithTax(0.2))

    const numbers = [1,2,3,4,5,6,7,8,9,10];
    const oddNumbers = numbers.filter(
      (num) => {
        return num % 2 === 1;
      }
    );
    console.log(oddNumbers);
    const evenNumbers = numbers.filter(num => num % 2 == 0 );
    console.log(evenNumbers);
  }


  someMethod() {
    let someValue: any;
    let anotherNumber: number;
    const aThirdNumber: number = 77;

    someValue = 12;
    anotherNumber = 44;
    anotherNumber = 22;
    someValue = "33";
  }

  exploringArrays() {
    // myArray1 and myArray2 are identical
    let myArray1: Array<any>;
    const myArray2: any[] = [1,2, "3", null, true];

    myArray1 = new Array<any>(5);

    console.log(myArray1);
    console.log(myArray2);

    // get item at index 1 from array
    console.log(myArray2[1]);

    // get slice between index 1 and 2 of the array
    console.log(myArray2.slice(1,2))

    // delete an item at index 1 of the array
    console.log(myArray2.splice(1, 1));
    console.log(myArray2);

    // add an item to the array (append to the end)
    console.log(myArray2.push(22));
    console.log(myArray2);

    // delete and return the last item from the array
    console.log(myArray2.pop());
    console.log(myArray2);

    // basic for loop (not idiomatic)
    for (let i = 0; i < 10; i++) {
      console.log(myArray2[i]);
    }

    // for each loop (idiomatic TypeScript way)
    // in = is the index
    // of = the value
    for (const v of myArray2) {
      console.log(v);
    }

    // while style loop
    let num = 0;
    while (num < 5) {
      console.log(num);
      num++;
    }

    // conditionals
    if (num < 5) {
      console.log('the number is less than 5');
    }
    else if (num == -1) {
      console.log('the number is invalid');
    }
    else {
      console.log('the number is 5 or greater');
    }
  }
}
