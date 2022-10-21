import {Component, Input, OnInit} from '@angular/core';
import {Layout, LayoutCapacity, Room} from "../../../model/Room";
import {FormBuilder, FormControl, FormGroup, Validators} from "@angular/forms";
import {DataService} from "../../../data.service";
import {Router} from "@angular/router";

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

  roomForm!: FormGroup;

  constructor(private formBuilder: FormBuilder, private dataService: DataService, private router: Router) {}

  ngOnInit(): void {
    this.roomForm = this.formBuilder.group(
      {
        roomName: [this.room.name,Validators.required],
        location: [this.room.location,[Validators.required,Validators.minLength(2)]],
      },
    );

    for (const layout of this.layouts) {
      const layoutCapacity = this.room.capacities.find((lc)=>{
        // @ts-ignore
        return lc.layout === Layout[layout]
      });
      const initialCapacity = layoutCapacity == null ? 0 : layoutCapacity.capacity;
      // @ts-ignore
      this.roomForm.addControl(`layout${layout}`, this.formBuilder.control(initialCapacity));
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

    // adding a new room
    if (this.room.id == null) {
      this.dataService.addNewRoom(this.room).subscribe(
        next => {
          this.router.navigate(['admin','rooms'], {queryParams:{action:'view', id:next.id}});
        }
      );
    } else {
      // updating an existing room
      this.dataService.updateRoom(this.room).subscribe(
        next =>{
          this.router.navigate(['admin','rooms'], {queryParams:{action:'view', id:next.id}});
        }
      );
    }
    console.log(this.room);
  }

  toTitle(s :string) :string {
      return s.split(' ').map(w => w[0].toUpperCase() + w.substr(1).toLowerCase()).join(' ');
  }

}
