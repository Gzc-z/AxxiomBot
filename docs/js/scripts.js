const images = document.querySelectorAll('.hero-carousel img');

let currentImage = 0;
let isAnimating = false;

function showImage(newIndex, direction) {
    if (isAnimating || newIndex === currentImage) return;

    isAnimating = true;

    const current = images[currentImage];
    const next = images[newIndex];

    next.style.transition = 'none';
    next.style.transform =
        direction === 'next'
            ? 'translateX(100%)'
            : 'translateX(-100%)';

    next.offsetHeight;

    next.style.transition = 'transform .5s ease';
    current.style.transition = 'transform .5s ease';

    current.style.transform =
        direction === 'next'
            ? 'translateX(-100%)'
            : 'translateX(100%)';

    next.style.transform = 'translateX(0)';

    setTimeout(() => {
        current.classList.remove('active');

        // reset past image
        current.style.transition = 'none';
        current.style.transform =
            direction === 'next'
                ? 'translateX(100%)'
                : 'translateX(-100%)';

        next.classList.add('active');

        currentImage = newIndex;
        isAnimating = false;
    }, 500);
}

function nextImage() {
    showImage((currentImage + 1) % images.length, 'next');
}

function prevImage() {
    showImage(
        (currentImage - 1 + images.length) % images.length,
        'prev'
    );
}

