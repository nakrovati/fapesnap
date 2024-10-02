import type { ParentComponent } from "solid-js";

const DefaultLayout: ParentComponent = (props) => {
  return (
    <div class="container py-16">
      <main>{props.children}</main>
    </div>
  );
};

export default DefaultLayout;
