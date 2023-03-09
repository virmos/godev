const slider = () => {
    const slides = document.querySelectorAll(".custom-slide");
    const btnLeft = document.querySelector(".custom-slider__btn--left");
    const btnRight = document.querySelector(
        ".custom-slider__btn--right"
    );
    const dotContainer = document.querySelector(".dots");

    let curSlide = 0;
    const maxSlide = slides.length;

    // Functions
    const activateDot = function (slide) {
        document
            .querySelectorAll(".dots__dot")
            .forEach((dot) =>
                dot.classList.remove("dots__dot--active")
            );

        document
            .querySelector(`.dots__dot[data-slide="${slide}"]`)
            .classList.add("dots__dot--active");
    };

    const goToSlide = function (slide) {
        slides.forEach(
            (s, i) =>
                (s.style.transform = `translateX(${
                    100 * (i - slide)
                }%)`)
        );
    };

    // Next slide
    const nextSlide = function () {
        if (curSlide === maxSlide - 1) {
            curSlide = 0;
        } else {
            curSlide++;
        }

        goToSlide(curSlide);
        activateDot(curSlide);
    };

    const prevSlide = function () {
        if (curSlide === 0) {
            curSlide = maxSlide - 1;
        } else {
            curSlide--;
        }
        goToSlide(curSlide);
        activateDot(curSlide);
    };

    const init = function () {
        goToSlide(0);
        activateDot(0);
    };
    init();

    // Event handlers
    btnRight.addEventListener("click", nextSlide);
    btnLeft.addEventListener("click", prevSlide);

    document.addEventListener("keydown", function (e) {
        if (e.key === "ArrowLeft") prevSlide();
        e.key === "ArrowRight" && nextSlide();
    });

    dotContainer.addEventListener("click", function (e) {
        if (e.target.classList.contains("dots__dot")) {
            const { slide } = e.target.dataset;
            goToSlide(slide);
            activateDot(slide);
        }
    });
};
export default slider;
