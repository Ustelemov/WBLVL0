interface OrderOut {
    order_uid: string;
    entry: string;
    total_price: number;
    customer_id: string;
    track_number: string;
    delivery_service: string;
}

const getOrderUrl = "http://localhost:8000/orders/order?uuid="
const getOrdersUrl = "http://localhost:8000/orders"
const uuid_input = document.getElementsByClassName('searchInput')[0] as HTMLInputElement
const img_404 = document.getElementsByClassName('imgbox')[0] as HTMLElement
const table = document.getElementsByClassName('tablebox')[0] as HTMLElement

img_404.style.display = "none"
table.style.display = "none"


const order_uid = document.getElementById('order_uid') as HTMLTableCellElement
const entry = document.getElementById('entry') as HTMLTableCellElement
const total_price = document.getElementById('total_price') as HTMLTableCellElement
const custumer_id = document.getElementById('customer_id') as HTMLTableCellElement
const delivery_service = document.getElementById('delivery_service') as HTMLTableCellElement
const track_number = document.getElementById('track_number') as HTMLTableCellElement


function getOrder<T>(url: string): Promise<T> {

    let headers = new Headers()
    headers.append('Content-Type', 'application/json');
    headers.append('Accept', 'application/json');

    let full_url = url

    return fetch(full_url, {
        credentials: 'include',
        method: 'GET',
        headers: headers
      })
      .then(response => {
        if (!response.ok) {
          throw new Error(response.statusText)
        }
        return response.json() as Promise<T>
      })
      .then(response => {
          return response
      })
}



function onClick() {
    
    let uuid = getOrderUrl + uuid_input.value

    getOrder(uuid)
      .then((response:OrderOut) => {
      
        order_uid.innerHTML = response.order_uid
        custumer_id.innerHTML = response.customer_id
        track_number.innerHTML = response.track_number
        entry.innerHTML = response.entry
        total_price.innerHTML = String(response.total_price)
        delivery_service.innerHTML = response.delivery_service


        img_404.style.display = "none"
        table.style.display = "block"


      })
      .catch(error => {
        console.log(error)
        img_404.style.display = "grid"
        table.style.display = "none"
      })
}

function onClickRandom() {
    
  let uuid = getOrdersUrl

  getOrder(uuid)
    .then((response:string[]) => {
      img_404.style.display = "none"
      table.style.display = "none"
      let random = Math.floor(Math.random() * response.length);

      if (response.length==0) {
        uuid_input.value = "No uuid is found"
      } else {
        uuid_input.value = response[random]
      }


    })
    .catch(error => {
      console.log(error)
      img_404.style.display = "none"
      table.style.display = "none"    
    })
}


