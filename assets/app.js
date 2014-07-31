var Map;

Map = (function() {
  function Map(elID) {
    var defaultCenter = new google.maps.LatLng(-23.5472, -46.6344),
        options = {
          zoom: 11,
          center: defaultCenter
        };

    this.gMap        = new google.maps.Map(document.getElementById(elID), options);

    this.bounds      = new google.maps.LatLngBounds();
    this.markers     = {};
    this.infoContent = {};
    this.infoWindow  = new google.maps.InfoWindow();
    this.oms         = new OverlappingMarkerSpiderfier(this.gMap);

    _this = this;
    this.oms.addListener('click', function(marker, event) {
      _this.infoWindow.setContent(marker.desc);
      _this.infoWindow.open(_this.gMap, marker);
    });
  }

  Map.prototype.extendBounds = function(latLng) {
    this.bounds.extend(latLng);
    this.gMap.panToBounds(this.bounds);
    this.gMap.fitBounds(this.bounds);
  };

  Map.prototype.setContent = function(index, data) {
    this.infoContent[index] = data
  }

  Map.prototype.addMarker = function(index, marker) {
    this.markers[index] = marker
    this.markers[index].setMap(this.gMap);
    this.oms.addMarker(this.markers[index]);
    _this = this;

    google.maps.event.addListener(this.markers[index], 'click', function(key) {
      return function(){
         _this.infoWindow.setContent(_this.infoContent[key]);
         _this.infoWindow.open(_this.gMap, _this.markers[key]);
      }
    }(index));
  };

  return Map;
})();

function initialize() {
  window.map = new Map("map-canvas");
  getPropertiesTo(window.location.pathname);
}

function getPropertiesTo(path) {
  oboe({
      url: "/api/imoveis"+window.location.pathname
  }).node(".*", function(){
  }).done(function(property) {
    if (property.lat != 0){
      var propertyLatLng = new google.maps.LatLng(property.lat, property.lng),
          index = property.id;

      window.map.extendBounds(propertyLatLng);
      window.map.setContent(index, '<div id="content">'+
         '<div id="siteNotice">'+
         '</div>'+
         '<h1 id="firstHeading" class="firstHeading">'+property.address+'</h1>'+
         '<div id="bodyContent">'+
           JSON.stringify(property, undefined, 2)+
         '</div>'+
         '</div>');

      window.map.addMarker(index, new google.maps.Marker({
        position: propertyLatLng,
        title: property.address
      }))
    }
  }).fail(function(err){
    console.log("E", err.thrown)
  })
}

  // function centralize() {
  //     window.setTimeout(function() {
  //       radius = window.location.hash.replace("#", "")
  //       if (radius.length < 1)  { radius = 2.0 }

  //       center = map.getCenter();
  //       getImoveis(center.lat(), center.lng(), radius);

  //       var request = {
  //         location: center,
  //         radius: radius*1000
  //                 };

  //       service = new google.maps.places.PlacesService(map);
  //       service.nearbySearch(request, callback);
  //     }, 2000);
  // }


// function callback(results, status) {
//   if (status == google.maps.places.PlacesServiceStatus.OK) {
//     for (var i = 0; i < results.length; i++) {
//       var place = results[i];

//       var gpmarker    = new google.maps.MarkerImage(place.icon, null, null, null, new google.maps.Size(25, 25));

//         var marker = new google.maps.Marker({
//               position: place.geometry.location,
//               title: place.name,
//               icon: gpmarker
//             });


//             marker.setMap(map);
//             window.markers.push(marker);

//     console.log(place)

//     }
//   }
// }




google.maps.event.addDomListener(window, 'load', initialize);
