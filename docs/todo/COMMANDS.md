# Commands

The system supports a variety of commands that can be sent from the client to the API to control the presentation and manage slides. All commands are documented in this file, which provides detailed information on the available commands, their parameters, and expected responses and broadcasts events.

Unless otherwise specified, all commands are sent as JSON objects through a WebSocket connection to the API. Each command should include a `command` field specifying the action to be performed, along with any necessary parameters.

## Command Structure

```json
{
    "command": "go_to_slide", // The command to be executed, e.g., "next_slide", "previous_slide", "annotate", etc.
    "parameters": { // An object containing any parameters required for the command.
        "slide_number": 3
    }
}
```

## Initialization Commands

- [x] **init_presentation**: Initializes a new presentation session.
  - Parameters:
    - `presentation_id`: A unique identifier for the presentation room.
  - Response: `{ "data": { "presentation_id": <presentation_id>, "current_index": 0, "slides": [<slide_object>] } }`
  - Broadcast: `{ "event": "presentation_initialized", "data": { "presentation_id": <presentation_id>, "current_index": 0, "slides": [<slide_object>] } }`
    - Slide Object Structure:
      ```json
      [
          {
              "id": "string",
              "content": "string"
          }
      ]
      ```

- [x] **join_room**: Subscribes to a presentation room to receive broadcast events.
  - Parameters:
    - `presentation_id`: The unique identifier of the presentation room to join.
  - Response: `{ "data": { "presentation_id": <presentation_id>, "current_index": <slide_number>, "slides": [<slide_object>] } }`
  - Broadcast: None (the server starts forwarding existing broadcasts to this connection).

## Slide Control Commands

- [x] **next_slide**: Advances to the next slide in the presentation.
  - Parameters: None
  - Broadcast: `{ "event": "slide_changed", "data": { "current_slide": <slide_number> } }`
- [x] **prev_slide**: Returns to the previous slide in the presentation.
  - Parameters: None
  - Broadcast: `{ "event": "slide_changed", "data": { "current_slide": <slide_number> } }`
- [x] **go_to_slide**: Jumps to a specific slide in the presentation.
  - Parameters:
    - `slide_number`: The number of the slide to jump to.
  - Broadcast: `{ "event": "slide_changed", "data": { "current_slide": <slide_number> } }`
- [ ] **hide_slide**: Hides the current slide from the presentation view.
  - Parameters: None
  - Response: `{ "data": { "current_slide": <slide_number>, "hidden": true } }`
  - Broadcast: `{ "event": "slide_hidden", "data": { "current_slide": <slide_number> } }`
- [ ] **show_slide**: Shows the current slide in the presentation view.
  - Parameters: None
  - Response: `{ "data": { "current_slide": <slide_number>, "hidden": false } }`
  - Broadcast: `{ "event": "slide_shown", "data": { "current_slide": <slide_number> } }`

## Annotation Commands

All coordinates are expressed as percentages (0-100) relative to the slide dimensions.

- [ ] **annotate**: Adds an annotation to the current slide.
  - Parameters:
    - `id`: A unique identifier for the annotation.
    - `points`: An array of points representing the annotation, where each point is an object with `x` and `y` coordinates.
    - `color`: The color of the annotation (e.g., "#FF0000" for red).
    - `thickness`: The thickness of the annotation line.
  - Response: `{ "data": { "current_slide": <slide_number>, "annotations": [<annotation_points>] } }`
  - Broadcast: `{ "event": "annotation_added", "data": { "current_slide": <slide_number>, "annotations": [<annotation_points>] } }`
- [ ] **clear_annotations**: Clears all annotations from the current slide.
  - Parameters: None
  - Response: `{ "data": { "current_slide": <slide_number>, "annotations": [] } }`
  - Broadcast: `{ "event": "annotations_cleared", "data": { "current_slide": <slide_number> } }`
- [ ] **undo_annotation**: Undoes the last annotation made on the current slide.
  - Parameters: None
  - Response: `{ "data": { "current_slide": <slide_number>, "annotations": [<remaining_annotations_points>] } }`
  - Broadcast: `{ "event": "annotation_undone", "data": { "current_slide": <slide_number>, "annotations": [<remaining_annotations_points>] } }`
- [ ] **redo_annotation**: Redoes the last undone annotation on the current slide.
  - Parameters: None
  - Response: `{ "data": { "current_slide": <slide_number>, "annotations": [<updated_annotations_points>] } }`
  - Broadcast: `{ "event": "annotation_redone", "data": { "current_slide": <slide_number>, "annotations": [<updated_annotations_points>] } }`

## Spin Wheel Commands

- [ ] **spin_wheel**: Spins the interactive wheel to randomly select a student.
  - Parameters: None
  - Response: `{ "data": { "selected_student": <student_id> } }`
  - Broadcast: `{ "event": "wheel_spun", "data": { "selected_student": <student_id> } }`
- [ ] **target_student**: Selects a specific student while maintaining the element of surprise for the rest of the class.
  - Parameters:
    - `student_id`: The ID of the student to be selected.
  - Response: `{ "data": { "selected_student": <student_id> } }`
  - Broadcast: `{ "event": "student_targeted", "data": { "selected_student": <student_id> } }`
- [ ] **show_wheel**: Displays the spin wheel on the presentation view.
  - Parameters: None
  - Response: `{ "data": { "wheel_visible": true } }`
  - Broadcast: `{ "event": "wheel_shown" }`
- [ ] **hide_wheel**: Hides the spin wheel from the presentation view.
  - Parameters: None
  - Response: `{ "data": { "wheel_visible": false } }`
  - Broadcast: `{ "event": "wheel_hidden" }`
- [ ] **add_student_to_wheel**: Adds a student to the spin wheel for selection.
  - Parameters:
    - `student_id`: The ID of the student to be added to the wheel.
  - Response: `{ "data": { "wheel_students": [<updated_student_list>] } }`
  - Broadcast: `{ "event": "student_added_to_wheel", "data": { "wheel_students": [<updated_student_list>] } }`
- [ ] **remove_student_from_wheel**: Removes a student from the spin wheel.
  - Parameters:
    - `student_id`: The ID of the student to be removed from the wheel.
  - Response: `{ "data": { "wheel_students": [<updated_student_list>] } }`
  - Broadcast: `{ "event": "student_removed_from_wheel", "data": { "wheel_students": [<updated_student_list>] } }`

