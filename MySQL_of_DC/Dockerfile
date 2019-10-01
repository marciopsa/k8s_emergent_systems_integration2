
FROM ubuntu:18.10

RUN apt update

RUN apt install mysql-server -y

RUN apt install libmysqlclient20

#CMD chown -R mysql:mysql /var/lib/mysql
#RUN /etc/init.d/mysql start

CMD chown -R mysql:mysql /var/lib/mysql
CMD usermod -d /var/lib/mysql/ mysql
RUN /etc/init.d/mysql start

CMD service mysql stop
CMD service mysql start





COPY ./dana ./home/marciopsa/dana 
COPY ./interscity_dana_project ./home/marciopsa/interscity_dana_project
COPY start.sh /home/marciopsa/interscity_dana_project/pal/start.sh
COPY assembly.sh /home/marciopsa/interscity_dana_project/pal/assembly.sh
ENV DANA_HOME=/home/marciopsa/dana/
ENV PATH="${PATH}:/home/marciopsa/dana/"
WORKDIR /home/marciopsa/interscity_dana_project/pal/



#ENTRYPOINT ["/bin/bash", "./assembly.sh"]
ENTRYPOINT ["/bin/bash", "./assembly.sh"]

#RUN mysqld
