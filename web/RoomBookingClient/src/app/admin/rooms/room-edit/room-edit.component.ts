import {Component, EventEmitter, Input, OnDestroy, OnInit, Output} from '@angular/core';
import {Layout, LayoutCapacity, Room} from "../../../model/Room";
import {FormBuilder, FormControl, FormGroup, Validators} from "@angular/forms";
import {DataService} from "../../../data.service";
import {Router} from "@angular/router";
import {FormResetService} from "../../../form-reset.service";
import {Subscription} from "rxjs";

@Component({
  selector: 'app-room-edit',
  templateUrl: './room-edit.component.html',
  styleUrls: ['./room-edit.component.css']
})
export class RoomEditComponent implements OnInit, OnDestroy {

  @Input()
  room!: Room;

  @Output()
  dataChanged = new EventEmitter()

  layouts = Object.keys(Layout);
  layoutEnum = Object(Layout);

  roomForm!: FormGroup;

  roomFormReset !: Subscription;

  message = "Plase wait..."

  constructor(private formBuilder: FormBuilder,
              private dataService: DataService,
              private router: Router,
              private formResetService: FormResetService) {}

  ngOnInit(): void {
   this.initializeForm();
   // this is used to clear the form
   this.roomFormReset = this.formResetService.resetRoomFormEvent.subscribe(
     room => {
       this.room = room;
     this.initializeForm();
   });
  }
  ngOnDestroy() {
    this.roomFormReset.unsubscribe();
  }

  initializeForm() {
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

  onCancel() {
    this.router.navigate(['admin','rooms'], {queryParams:{id: this.room.id, action:'view'}});
  }

  onSubmit() {
    this.message = "Saving..."
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
          this.dataChanged.emit();
          this.router.navigate(['admin','rooms'], {queryParams:{action:'view', id:next.id}});
        },
        error => {
          this.message = 'Something went wrong, please try again';
        }
      );
    } else {
      // updating an existing room
      this.dataService.updateRoom(this.room).subscribe(
        next =>{
          this.dataChanged.emit();
          this.router.navigate(['admin','rooms'], {queryParams:{action:'view', id:next.id}});
        },
        error => {
          this.message = 'Something went wrong, please try again';
        }
      );
    }
  }

  toTitle(s :string) :string {
      return s.split(' ').map(w => w[0].toUpperCase() + w.substr(1).toLowerCase()).join(' ');
  }

}
