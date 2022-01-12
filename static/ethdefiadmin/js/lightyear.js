var pageLoader = function($mode) {
	var $loadingEl = parent.jQuery('#lyear-loading');
	    $mode      = $mode || 'show';
	if ($mode === 'show') {
		if ($loadingEl.length) {
			$loadingEl.fadeIn(250);
		} else {
			parent.jQuery('body').prepend('<div id="lyear-loading"><div class="spinner-border text-primary" role="status"><span class="sr-only">Loading...</span></div></div>');
		}
	} else if ($mode === 'hide') {
		if ($loadingEl.length) {
			$loadingEl.fadeOut(250);
		}
	}
	return false;
};
