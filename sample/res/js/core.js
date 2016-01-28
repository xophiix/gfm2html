/****************************************
  SIDEBAR
****************************************/
$(window).load(function(){

  function loop(links, ul) {
    var html = '';
    $.each(links, function(key,obj){

      var $li = $('<li>');
      $li.appendTo(ul);

      var $page = location.href.split('/');
      var $class = ($page[$page.length - 1].split('.html')[0] == obj['Link']) ? 'current' : '';

      if(obj['Link'] == 'null'){
        $li.text(obj['Title']).addClass('nl').wrapInner('<span></span>');
      }
      else {
        var $a = $('<a>').attr('href',obj['Link']).text(obj['Title']);
        $($class != '') ? $a.attr('id',$class).attr('class',$class) : '';
        $a.appendTo($li);
      }
      if(obj['Children'] && obj['Children'].length > 0) {
        $li.prepend('<div class="arrow"></div>');
        var ul2 = $('<ul>').appendTo($li);
        loop(obj['Children'], ul2);
      }
    });
  }

  var data = toc;
  var ul = $('<ul/>');
  loop(data['Children'], ul);
  $('.sidebar-menu .toc').append(ul);

  // Prepare sidebar list
  $('.sidebar-menu .toc').find('li:has(ul)').children('ul').hide();
  $('.sidebar-menu .toc').find('.arrow').addClass('collapsed');

  // If no current, find parent
  if($('.sidebar-menu .toc .current').size() == 0){

    var url = location.href.split('/'),
        page = url[url.length - 1].split('.html')[0];

    // Find last occurrence of either '-' or '.' and split on this
    var splitAtPos = page.lastIndexOf('-');
    splitAtPos = Math.max(splitAtPos, page.lastIndexOf('.'));

    var parentPageName = page;
    if (splitAtPos != -1)
        parentPageName = page.substring(0, splitAtPos);
    $('.sidebar-menu .toc a[href="'+ parentPageName +'.html"]').addClass('current');
  }

  // Auto expand current section in sidebar
  var parents = $('.sidebar-menu .toc .current').parents('ul');
  if(parents.length > 0){
    parents.addClass('exp');
    $('ul.exp').removeAttr('style');
    $('.sidebar-menu ul.exp:first').removeAttr('class');
    $('ul.exp').parent().find('.arrow:first').removeClass('collapsed').addClass('expanded');
  }

  // Expand or collapse
  $('.sidebar-menu .toc li:has(ul) .arrow').click(function(e){
    var arrow = $(this);
    if(arrow.hasClass('expanded')){
      arrow.removeClass('expanded');
      arrow.addClass('collapsed');
    }
    else {
      arrow.removeClass('collapsed');
      arrow.addClass('expanded');
    }
    arrow.parent().find('ul:first').toggle();
  });

  // For items with no link
  $('.sidebar-menu .toc .nl span').click(function(e){
    var arrow = $(this).prev('.arrow');
    if(arrow.hasClass('expanded')){
      arrow.removeClass('expanded');
      arrow.addClass('collapsed');
    }
    else {
      arrow.removeClass('collapsed');
      arrow.addClass('expanded');
    }
    arrow.parent().find('ul:first').toggle();
  });
});