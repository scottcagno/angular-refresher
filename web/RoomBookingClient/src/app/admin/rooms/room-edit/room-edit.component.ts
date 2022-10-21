import {Component, Input, OnInit} from '@angular/core';
import {Layout, LayoutCapacity, Room} from "../../../model/Room";
import {FormControl, FormGroup} from "@angular/forms";

@Component({
  selector: 'app-room-edit',
  templateUrl: './room-edit.component.html',
  styleUrls: ['./room-edit.component.css']
})
export class RoomEditComponent implements OnInit {

  @Input()
  room!: Room;

  layouts = Object.keys(Layout);
  layoutEnum = Object(Layout);

  roomForm = new FormGroup(
    {
      roomName : new FormControl('name'),
      location : new FormControl('location'),
    }
  );

  constructor() {}

  ngOnInit(): void {
    this.roomForm.patchValue({
      roomName:this.room.name,
      location: this.room.location,
    });
    for (const layout of this.layouts) {
      // @ts-ignore
      this.roomForm.addControl(`layout${layout}`, new FormControl(`layout${layout}`));
    }
  }

  onCancel() {}

  onSubmit() {
    this.room.name = <string>this.roomForm.controls['roomName'].value;
    this.room.location = <string>this.roomForm.controls['location'].value;
    this.room.capacities = new Array<LayoutCapacity>();
    for (const layout of this.layouts) {
      const layoutCapacity = new LayoutCapacity();
      // @ts-ignore
      layoutCapacity.layout = Layout[layout];
      // @ts-ignore
      layoutCapacity.capacity = <number>this.roomForm.controls[`layout${layout}`].value;
      this.room.capacities.push(layoutCapacity);
    }
    console.log(this.room);
  }

  toTitle(s :string) :string {
      return s.split(' ').map(w => w[0].toUpperCase() + w.substr(1).toLowerCase()).join(' ');
  }

}
