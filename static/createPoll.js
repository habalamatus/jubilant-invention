$(document).ready(function () {
  var max_fields = 10;
  var wrapper = $(".container1");
  var add_button = $(".add_form_field");

  var x = 1;
  $(add_button).click(function (e) {
    e.preventDefault();
    if (x < max_fields) {
      x++;
      $(wrapper).append(
        '<div class=" flex"><input type="text" name="option" required class="mb-3 mr-3 bg-gray-200 rounded w-full text-gray-700 focus:outline-none border-b-4 border-gray-300 focus:border-blue-600 transition duration-500 px-3 pb-3"><div class="m-auto"><a href="#" class="delete ">Remove</a></div>'
      ); //add input box
    } else {
      alert("You Reached the limits");
    }
  });

  $(wrapper).on("click", ".delete", function (e) {
    e.preventDefault();
    $(this).parent("div").parent("div").remove();
    x--;
  });
});

const tomorrow = new Date();
tomorrow.setDate(tomorrow.getDate() + 30);
document.getElementById("end").valueAsDate = tomorrow;
