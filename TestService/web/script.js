var getOrderUrl = "http://localhost:8000/orders/order?uuid=";
var getOrdersUrl = "http://localhost:8000/orders";
var uuid_input = document.getElementsByClassName('searchInput')[0];
var img_404 = document.getElementsByClassName('imgbox')[0];
var table = document.getElementsByClassName('tablebox')[0];
img_404.style.display = "none";
table.style.display = "none";
var order_uid = document.getElementById('order_uid');
var entry = document.getElementById('entry');
var total_price = document.getElementById('total_price');
var custumer_id = document.getElementById('customer_id');
var delivery_service = document.getElementById('delivery_service');
var track_number = document.getElementById('track_number');
function getOrder(url) {
    var headers = new Headers();
    headers.append('Content-Type', 'application/json');
    headers.append('Accept', 'application/json');
    var full_url = url;
    return fetch(full_url, {
        credentials: 'include',
        method: 'GET',
        headers: headers
    })
        .then(function (response) {
        if (!response.ok) {
            throw new Error(response.statusText);
        }
        return response.json();
    })
        .then(function (response) {
        return response;
    });
}
function onClick() {
    var uuid = getOrderUrl + uuid_input.value;
    getOrder(uuid)
        .then(function (response) {
        order_uid.innerHTML = response.order_uid;
        custumer_id.innerHTML = response.customer_id;
        track_number.innerHTML = response.track_number;
        entry.innerHTML = response.entry;
        total_price.innerHTML = String(response.total_price);
        delivery_service.innerHTML = response.delivery_service;
        img_404.style.display = "none";
        table.style.display = "block";
    })["catch"](function (error) {
        console.log(error);
        img_404.style.display = "grid";
        table.style.display = "none";
    });
}
function onClickRandom() {
    var uuid = getOrdersUrl;
    getOrder(uuid)
        .then(function (response) {
        img_404.style.display = "none";
        table.style.display = "none";
        var random = Math.floor(Math.random() * response.length);
        if (response.length == 0) {
            uuid_input.value = "No uuid is found";
        }
        else {
            uuid_input.value = response[random];
        }
    })["catch"](function (error) {
        console.log(error);
        img_404.style.display = "none";
        table.style.display = "none";
    });
}
