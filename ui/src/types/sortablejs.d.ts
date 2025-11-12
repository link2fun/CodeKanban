declare module 'sortablejs' {
  export interface SortableEvent extends Event {
    oldIndex?: number;
    newIndex?: number;
    item?: HTMLElement;
    oldDraggableIndex?: number;
    newDraggableIndex?: number;
    [key: string]: any;
  }

  export default class Sortable {
    el: HTMLElement;
    constructor(element: HTMLElement, options?: any);
    option(name: string, value?: any): any;
    destroy(): void;
    static create(element: HTMLElement, options?: any): Sortable;
  }
}
