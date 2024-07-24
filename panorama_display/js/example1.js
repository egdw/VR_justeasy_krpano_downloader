window.onload = function() {
	document.getElementById('go').addEventListener('click', loadPredefinedPanorama, false);

	document.getElementById('pano').addEventListener('change', upload, false);
};

// Load the predefined panorama
function loadPredefinedPanorama(evt) {
	evt.preventDefault();

	var div = document.getElementById('container');

	var PSV = new PhotoSphereViewer({
		// Path to the panorama
		panorama: 'sun.jpg',

		// Container
		container: div,

		// Deactivate the animation
		time_anim: false,

		// Display the navigation bar
		navbar: true,

		// Resize the panorama
		size: {
			width: '100%',
			height: '500px'
		}
	});
}

// Load a panorama stored on the user's computer
function upload() {
	// Retrieve the chosen file and create the FileReader object
	var file = document.getElementById('pano').files[0];
	var reader = new FileReader();

	reader.onload = function() {
		var div = document.getElementById('your-pano');

		var PSV = new PhotoSphereViewer({
			// Panorama, given in base 64
			panorama: reader.result,

			// Container
			container: div,

			// Deactivate the animation
			time_anim: false,

			// Display the navigation bar
			navbar: true,

			// Resize the panorama
			size: {
				width: '100%',
				height: '500px'
			}
		});
	};

	reader.readAsDataURL(file);
}
