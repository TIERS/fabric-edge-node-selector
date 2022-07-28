FROM osrf/ros:foxy-desktop-focal

SHELL ["/bin/bash", "-c"]

WORKDIR /app

#COPY ros2_ws ros2_ws/
#RUN cd ros2_ws && \
RUN source /opt/ros/foxy/setup.bash
    #colcon build

COPY ros_entrypoint.sh /

ENTRYPOINT ["/ros_entrypoint.sh"]
CMD ["bash"]