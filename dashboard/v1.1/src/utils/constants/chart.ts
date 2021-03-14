export const JITTER_OPTIONS = {
  chart: {
    type: 'area',
    animations: {
      enabled: true,
      easing: 'easeinout',
      speed: 800,
      animateGradually: {
        enabled: true,
        delay: 150
      },
      dynamicAnimation: {
        enabled: true,
        speed: 350
      }
    },
    background: '#fff'
  },
  yaxis: {
    title: {
      text: 'milliseconds'
    }
  },
  xaxis: {
    title: {
      text: 'Time'
    }
  }
};

export const PING_OPTIONS = {
  chart: {
    type: 'area',
    xaxis: {
      type: 'category',
      categories: [],
      labels: {
        show: true,
        rotate: 0,
        rotateAlways: false,
        hideOverlappingLabels: true,
        trim: true
      }
    },
    animations: {
      enabled: true,
      easing: 'easeinout',
      speed: 800,
      animateGradually: {
        enabled: true,
        delay: 150
      },
      dynamicAnimation: {
        enabled: true,
        speed: 350
      }
    },
    background: '#fff'
  },
  yaxis: {
    title: {
      text: 'milliseconds'
    }
  },
  xaxis: {
    title: {
      text: 'Time'
    }
  }
};

export const DELAY_OPTIONS = {
  chart: {
    type: 'area',
    animations: {
      enabled: true,
      easing: 'easeinout',
      speed: 800,
      animateGradually: {
        enabled: true,
        delay: 150
      },
      dynamicAnimation: {
        enabled: true,
        speed: 350
      }
    },
    background: '#fff'
  },
  yaxis: {
    title: {
      text: 'Response (in ms)'
    }
  },
  xaxis: {
    title: {
      text: 'Time'
    }
  }
};

export const RESPONSE_LENGTH_OPTIONS = {
  chart: {
    type: 'area',
    animations: {
      enabled: true,
      easing: 'easeinout',
      speed: 800,
      animateGradually: {
        enabled: true,
        delay: 150
      },
      dynamicAnimation: {
        enabled: true,
        speed: 350
      }
    },
    background: '#fff'
  },
  yaxis: {
    title: {
      text: 'Length'
    }
  },
  xaxis: {
    title: {
      text: 'Time'
    }
  }
};
