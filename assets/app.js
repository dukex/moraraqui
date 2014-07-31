  function initialize() {
    var defaultMapCenter = new google.maps.LatLng(-23.5472, -46.6344),
    mapOptions = {
      zoom: 14,
      center: defaultMapCenter
    };

    window.bounds  = new google.maps.LatLngBounds();
    window.map     = new google.maps.Map(document.getElementById("map-canvas"), mapOptions);
    window.markers = {}
    window.infoWindows = {}
    window.infoContent = {}


    getDataTo(window.location.pathname);

    window.iw = new google.maps.InfoWindow();
    window.oms = new OverlappingMarkerSpiderfier(map);

    oms.addListener('click', function(marker, event) {
      iw.setContent(marker.desc);
      iw.open(map, marker);
    });
  }

  function getDataTo(path) {
    oboe({
      url: "/api/imoveis"+window.location.pathname
    }).node(".*", function(){
    }).done(function(imovel) {
      if (imovel.lat != 0){
        var imovelLatLng = new google.maps.LatLng(imovel.lat,imovel.lng),
            i = imovel.id;

        // var image = {
        //   url: "/house.png"
        // };

        window.bounds.extend(imovelLatLng);
        window.infoContent[i] = '<div id="content">'+
           '<div id="siteNotice">'+
           '</div>'+
           '<h1 id="firstHeading" class="firstHeading">'+imovel.address+'</h1>'+
           '<div id="bodyContent">'+
            JSON.stringify(imovel, undefined, 2);
           '</div>'+
           '</div>';

        window.markers[i] = new google.maps.Marker({
          position: imovelLatLng,
          title: imovel.address
          // icon: image
        });

        window.markers[i].setMap(map);
        window.oms.addMarker(window.markers[i]);



        google.maps.event.addListener(window.markers[i], 'click', function(infoKey) {
          return function(){
            window.iw.setContent(window.infoContent[infoKey]);
            window.iw.open(map, window.markers[i]);
          }
        }(i));
      }
    }).fail(function(err){
      console.log("E", err)
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