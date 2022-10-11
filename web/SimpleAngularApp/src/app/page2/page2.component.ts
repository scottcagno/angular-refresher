import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-page2',
  templateUrl: './page2.component.html',
  styleUrls: ['./page2.component.css']
})
export class Page2Component implements OnInit {

  hits!:number;

  constructor() {
    this.hits = 0;
  }

  ngOnInit(): void {
  }

  incrementHitCounter() {
    this.hits++;
  }

}
