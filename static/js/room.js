function Room(csrftoken, room_id){
    document.getElementById("check-availability").addEventListener("click", function () {
        let html = `
      <form action="/search-availability-json" method = "post" action="" id="check-availability-form" novalidate class="needs-validation">
        <input type="hidden" name="csrf_token" value="${csrftoken}">
        <input type="hidden" name="room_id" value="${room_id}">
  
          <div class="form-row">
              <div class="col">
                  <div id="reservation-date-modal">
                      <div class="row">
                          <div class="col">
                              <input type="text" disabled required class="form-control" name="start" id="start" placeholder="Arrival">
                          </div>
                          <div class="col">
                              <input type="text" disabled required class="form-control" name="end" id="end" placeholder="Departure">
                          </div>
                      </div>
                  </div>
              </div>
          </div>
      </form>`
      
        let attention = Prompt()
        attention.custom({
          html: html,
          title: "Pick your dates",
  
          willOpen: () => {
            const elem = document.getElementById('reservation-date-modal');
            const rangepicker = new DateRangePicker(elem, {
              format: "yyyy-mm-dd",
              minDate: new Date(),
              showOnFocus: true
            });
          },
  
          didOpen: ()=> {
                      document.getElementById('start').removeAttribute('disabled')
                      document.getElementById('end').removeAttribute('disabled')
                  },
  
          callback: function (result) {
  
            let form = document.getElementById("check-availability-form");
            let formData = new FormData(form);
  
            fetch('/search-availability-json', {
              method: "post",
              body: formData
            })
              .then(response => response.json())
              .then(data => {
                if (data.ok) {
                  attention.custom({
                    showConfirmButton : false,
                    showCancelButton : false,
                    icon: 'success',
                    title: `<p>Room is available</p><br>`,
                    html: `<p><a href="/book-room?`
                      +`id=` + data.room_id
                      +`&s=` + data.start_date
                      +`&e=` + data.end_date
                      +`" class="btn btn-primary"> Book Now </a></p> `
                  })
                }else {
                  attention.error({
                    title:`Room unavailable`,
                  })
                }
              })
          }
        })
  
      }) 
}
 