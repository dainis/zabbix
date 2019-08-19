#!/usr/bin/env ruby

# Common setup script for Vagrant and Travis CI.

ZABBIX_VERSION = ENV['ZABBIX_VERSION']
puts "Installing Zabbix #{ZABBIX_VERSION} ..."

def run(cmd)
  puts "# #{cmd}"
  res = system cmd
  if not res
    puts "#{res} #{$?}"
    exit res
  end
end

run('curl http://repo.zabbix.com/zabbix-official-repo.key 2>/dev/null| apt-key add -')
run('curl "https://keyserver.ubuntu.com/pks/lookup?op=get&search=0xb97b0afcaa1a47f044f244a07fcc7d46accc4cf8" 2>/dev/null| apt-key add -')
run("add-apt-repository 'deb http://repo.zabbix.com/zabbix/#{ZABBIX_VERSION}/ubuntu/ trusty main non-free contrib'")
run("LC_ALL=C.UTF-8 add-apt-repository -y ppa:ondrej/php")

# remove apt.postgresql.org repository and use only standard one,
# otherwise Zabbix 2.0 will not be installed with error:
#   The following packages have unmet dependencies:
#   zabbix-server-pgsql : Depends: libiodbc2 (>= 3.52.7) but it is not going to be installed
#   E: Unable to correct problems, you have held broken packages.
#run("add-apt-repository --remove 'deb http://apt.postgresql.org/pub/repos/apt/ precise-pgdg main'")
#run('apt-get purge -y postgresql.*')
#run('apt-get autoremove --purge')
#run('rm -fr /var/lib/postgresql /etc/postgresql')

run('apt-get update')
run('apt-get install -y postgresql-11')
run('apt-get install -y apache2 php7.3-fpm php7.3-mysql')
File.open('/etc/php/7.3/fpm/php.ini', 'a') do |f|
  f.puts '[Date]'
  f.puts 'date.timezone = UTC'
end
run('apt-get install -y zabbix-server-pgsql zabbix-frontend-php')
run('a2enmod proxy_fcgi setenvif')
run('a2enconf php7.3-fpm')

conf = File.read('/etc/dbconfig-common/zabbix-server-pgsql.conf')
password = /dbc_dbpass='(\w+)'/.match(conf)[1]

File.open('/usr/share/zabbix/conf/zabbix.conf.php', 'w') do |f| f.puts <<-END
  <?php
  // Zabbix GUI configuration file
  global $DB;

  $DB['TYPE']     = 'POSTGRESQL';
  $DB['SERVER']   = 'localhost';
  $DB['PORT']     = '0';
  $DB['DATABASE'] = 'zabbix';
  $DB['USER']     = 'zabbix';
  $DB['PASSWORD'] = '#{password}';

  // SCHEMA is relevant only for IBM_DB2 database
  $DB['SCHEMA'] = '';

  $ZBX_SERVER      = 'localhost';
  $ZBX_SERVER_PORT = '10051';
  $ZBX_SERVER_NAME = '';

  $IMAGE_FORMAT_DEFAULT = IMAGE_FORMAT_PNG;
  ?>
  END
end
